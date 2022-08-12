package destinationfetchersvc

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/config"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/kyma-incubator/compass/components/director/pkg/persistence"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"golang.org/x/sync/semaphore"

	"github.com/pkg/errors"
)

const (
	correlationIDPrefix = "sap.s4:communicationScenario:"
	regionLabelKey      = "region"
)

//go:generate mockery --name=UUIDService --output=automock --outpkg=automock --case=underscore --disable-version-string
// UUIDService missing godoc
type UUIDService interface {
	Generate() string
}

//go:generate mockery --name=DestinationRepo --output=automock --outpkg=automock --case=underscore --disable-version-string
// DestinationRepo missing godoc
type DestinationRepo interface {
	Upsert(ctx context.Context, in model.DestinationInput, id, tenantID, bundleID, revision string) error
	Delete(ctx context.Context, revision string) error
}

//go:generate mockery --name=LabelRepo --output=automock --outpkg=automock --case=underscore --disable-version-string
// LabelRepo missing godoc
type LabelRepo interface {
	GetSubdomainLabelForSubscribedRuntime(ctx context.Context, tenantID string) (*model.Label, error)
	GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID, key string) (*model.Label, error)
}

//go:generate mockery --name=BundleRepo --output=automock --outpkg=automock --case=underscore --disable-version-string
// BundleRepo missing godoc
type BundleRepo interface {
	GetBySystemAndCorrelationID(ctx context.Context, tenantID, systemName, systemURL, correlationID string) ([]*model.Bundle, error)
}

//go:generate mockery --name=TenantRepo --output=automock --outpkg=automock --case=underscore --disable-version-string
// TenantRepo missing godoc
type TenantRepo interface {
	GetBySubscribedRuntimes(ctx context.Context) ([]*model.BusinessTenantMapping, error)
}

// DestinationService missing godoc
type DestinationService struct {
	Transactioner      persistence.Transactioner
	UUIDSvc            UUIDService
	Repo               DestinationRepo
	BundleRepo         BundleRepo
	LabelRepo          LabelRepo
	DestinationsConfig config.DestinationsConfig
	APIConfig          DestinationServiceAPIConfig
	TenantRepo         TenantRepo
}

// GetSubscribedTenantIDs returns subscribed tenants
func (d *DestinationService) GetSubscribedTenantIDs(ctx context.Context) ([]string, error) {
	tenants, err := d.getSubscribedTenants(ctx)
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to get subscribed tenants")
		return nil, err
	}
	tenantIDs := make([]string, 0, len(tenants))
	for _, tenant := range tenants {
		tenantIDs = append(tenantIDs, tenant.ExternalTenant)
	}
	return tenantIDs, nil
}

func (d *DestinationService) getSubscribedTenants(ctx context.Context) ([]*model.BusinessTenantMapping, error) {
	tx, err := d.Transactioner.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to begin db transaction")
		return nil, err
	}
	defer d.Transactioner.RollbackUnlessCommitted(ctx, tx)
	ctx = persistence.SaveToContext(ctx, tx)

	tenants, err := d.TenantRepo.GetBySubscribedRuntimes(ctx)
	if err != nil {
		log.C(ctx).WithError(err).Error("An error occurred while getting subscribed tenants")
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		log.C(ctx).WithError(err).Error("An error occurred while committing database transaction")
		return nil, err
	}
	return tenants, nil
}

// SyncTenantDestinations syncs destinations for a given tenant
func (d *DestinationService) SyncTenantDestinations(ctx context.Context, tenantID string) error {
	subdomainLabel, err := d.getSubscribedSubdomainLabel(ctx, tenantID)
	if err != nil {
		return err
	}

	regionLabel, err := d.getRegionLabel(ctx, *subdomainLabel.Tenant)
	if err != nil {
		return err
	}

	subdomain := subdomainLabel.Value.(string)
	region := regionLabel.Value.(string)

	instanceConfig, ok := d.DestinationsConfig.RegionToInstanceConfig[region]
	if !ok {
		log.C(ctx).Errorf("No destination instance credentials found for region '%s'", region)
		return errors.New(fmt.Sprintf("No destination instance credentials found for region '%s'", region))
	}

	client, err := NewClient(instanceConfig, d.APIConfig, d.DestinationsConfig.OauthTokenPath, subdomain)
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to create Destination API client")
		return err
	}

	err = d.walkthroughPages(ctx, client, func(destinations []model.DestinationInput) error {
		log.C(ctx).Infof("Found %d destinations in tenant '%s'", len(destinations), tenantID)
		return d.mapDestinationsToTenant(ctx, *subdomainLabel.Tenant, destinations)
	})
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to sync destinations for tenant '%s'", tenantID)
		return err
	}

	return nil
}

func (d *DestinationService) mapDestinationsToTenant(ctx context.Context, tenant string, destinations []model.DestinationInput) error {
	tx, err := d.Transactioner.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to begin db transaction")
		return err
	}
	ctx = persistence.SaveToContext(ctx, tx)
	defer d.Transactioner.RollbackUnlessCommitted(ctx, tx)

	for _, destination := range destinations {
		correlationID := correlationIDPrefix + destination.CommunicationScenarioID
		bundles, err := d.BundleRepo.GetBySystemAndCorrelationID(ctx, tenant, destination.XFSystemName, destination.URL, correlationID)

		if len(bundles) == 0 {
			log.C(ctx).Infof("No bundles found for system '%s', url '%s', correlation id '%s'", destination.XFSystemName, destination.URL, correlationID)
			continue
		}

		if err != nil {
			log.C(ctx).WithError(err).Errorf("Failed to fetch bundle for system '%s', url '%s', correlation id '%s', tenant id '%s'",
				destination.XFSystemName, destination.URL, correlationID, tenant)
			continue
		}

		for _, bundle := range bundles {
			id := d.UUIDSvc.Generate()
			revision := d.UUIDSvc.Generate()
			if err := d.Repo.Upsert(ctx, destination, id, tenant, bundle.ID, revision); err != nil {
				log.C(ctx).WithError(err).Errorf(
					"Failed to insert destination with name '%s' for bunlde '%s' and tenant '%s' to DB",
					destination.Name, bundle.ID, tenant)
				continue
			}
		}
	}

	if err = tx.Commit(); err != nil {
		log.C(ctx).WithError(err).Error("Failed to commit database transaction")
		return err
	}
	return nil
}

type processFunc func([]model.DestinationInput) error

func (d *DestinationService) walkthroughPages(ctx context.Context, client *Client, process processFunc) error {
	hasMorePages := true

	for page := 1; hasMorePages; page++ {
		pageString := strconv.Itoa(page)
		resp, err := client.FetchTenantDestinationsPage(ctx, pageString)
		if err != nil {
			return errors.Wrap(err, "failed to fetch destinations page")
		}

		if err := process(resp.destinations); err != nil {
			return errors.Wrap(err, "failed to process destinations page")
		}

		hasMorePages = pageString != resp.pageCount
	}

	return nil
}

// FetchDestinationsSensitiveData returns sensitive data of destinations for a given tenant
func (d *DestinationService) FetchDestinationsSensitiveData(ctx context.Context, tenantID string, destinationNames []string) ([]byte, error) {
	subdomainLabel, err := d.getSubscribedSubdomainLabel(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	regionLabel, err := d.getRegionLabel(ctx, *subdomainLabel.Tenant)
	if err != nil {
		return nil, err
	}

	subdomain := subdomainLabel.Value.(string)
	region := regionLabel.Value.(string)
	log.C(ctx).Infof("Fetching data for subdomain: %s \n", subdomain)

	instanceConfig, ok := d.DestinationsConfig.RegionToInstanceConfig[region]
	if !ok {
		log.C(ctx).Errorf("No destination instance credentials found for region '%s'", region)
		return nil, errors.New(fmt.Sprintf("No destination instance credentials found for region '%s'", region))
	}
	client, err := NewClient(instanceConfig, d.APIConfig, d.DestinationsConfig.OauthTokenPath, subdomain)
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to create Destination API client")
		return nil, err
	}

	nameCount := len(destinationNames)
	results := make([][]byte, nameCount)
	weighted := semaphore.NewWeighted(d.APIConfig.GoroutineLimit)
	resChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		for _, destination := range destinationNames {
			if err := weighted.Acquire(ctx, 1); err != nil {
				log.C(ctx).WithError(err).
					Errorf("Could not acquire semaphore. Destination %s will not be fetched", destination)
				return
			}
			go fetchDestination(ctx, destination, weighted, client, resChan, errChan)
		}
	}()

	for i := 0; i < nameCount; {
		select {
		case err := <-errChan:
			return nil, err
		case res := <-resChan:
			results[i] = res
			i++
		}
	}

	combinedInfoJSON := bytes.Join(results, []byte(","))
	combinedInfoJSON = append(combinedInfoJSON, '}', '}')

	return append([]byte(`{ "destinations": {`), combinedInfoJSON...), nil
}

func fetchDestination(ctx context.Context, dest string, weighted *semaphore.Weighted,
	client *Client, resChan chan []byte, errChan chan error) {
	log.C(ctx).Infof("Fetching data for destination: %s \n", dest)
	defer weighted.Release(1)
	result, err := client.FetchDestinationSensitiveData(ctx, dest)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to fetch data for destination %s", dest)
		errChan <- err
		return
	}

	result = append([]byte("\""+dest+"\":"), result...)

	resChan <- result
}

func (d *DestinationService) getSubscribedSubdomainLabel(ctx context.Context, tenantID string) (*model.Label, error) {
	tx, err := d.Transactioner.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to begin db transaction")
		return nil, err
	}
	ctx = persistence.SaveToContext(ctx, tx)
	defer d.Transactioner.RollbackUnlessCommitted(ctx, tx)

	label, err := d.LabelRepo.GetSubdomainLabelForSubscribedRuntime(ctx, tenantID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			log.C(ctx).Errorf("No subscribed subdomain found for tenant '%s'", tenantID)
			return nil, apperrors.NewNotFoundErrorWithMessage(resource.Label, tenantID, fmt.Sprintf("tenant %s not found", tenantID))
		}
		log.C(ctx).WithError(err).Errorf("Failed to get subdomain for tenant '%s' from db", tenantID)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.C(ctx).WithError(err).Error("Failed to commit database transaction")
		return nil, err
	}

	return label, nil
}

func (d *DestinationService) getRegionLabel(ctx context.Context, tenantID string) (*model.Label, error) {
	tx, err := d.Transactioner.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Error("Failed to begin db transaction")
		return nil, err
	}
	ctx = persistence.SaveToContext(ctx, tx)
	defer d.Transactioner.RollbackUnlessCommitted(ctx, tx)

	region, err := d.LabelRepo.GetByKey(ctx, tenantID, model.TenantLabelableObject, tenantID, regionLabelKey)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to fetch region for tenant '%s'", tenantID)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.C(ctx).WithError(err).Error("Failed to commit database transaction")
		return nil, err
	}

	return region, nil
}
