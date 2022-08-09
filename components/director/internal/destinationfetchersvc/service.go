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

type UUIDService interface {
	Generate() string
}

type DestinationRepo interface {
	Upsert(ctx context.Context, in model.DestinationInput, id, tenantID, bundleID, revision string) error
	Delete(ctx context.Context, revision string) error
}

type LabelRepo interface {
	GetSubdomainLabelForSubscribedRuntime(ctx context.Context, subaccountId string) (*model.Label, error)
	GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID, key string) (*model.Label, error)
}

type BundleRepo interface {
	GetBySystemAndCorrelationId(ctx context.Context, tenantId, systemName, systemURL, correlationId string) ([]*model.Bundle, error)
}

type TenantRepo interface {
	GetBySubscibredRuntimes(ctx context.Context) ([]*model.BusinessTenantMapping, error)
}

type DestinationService struct {
	transact           persistence.Transactioner
	uuidSvc            UUIDService
	repo               DestinationRepo
	bundleRepo         BundleRepo
	labelRepo          LabelRepo
	tenantRepo         TenantRepo
	destinationsConfig config.DestinationsConfig
	apiConfig          APIConfig
}

func NewDestinationService(transact persistence.Transactioner, uuidSvc UUIDService, destRepo DestinationRepo,
	bundleRepo BundleRepo, labelRepo LabelRepo, tenantRepo TenantRepo, destinationsConfig config.DestinationsConfig, apiConfig APIConfig) *DestinationService {
	return &DestinationService{
		transact:           transact,
		uuidSvc:            uuidSvc,
		repo:               destRepo,
		bundleRepo:         bundleRepo,
		labelRepo:          labelRepo,
		tenantRepo:         tenantRepo,
		destinationsConfig: destinationsConfig,
		apiConfig:          apiConfig,
	}
}

func (d DestinationService) SyncSubaccountDestinations(ctx context.Context, subaccountID string) error {
	subdomainLabel, err := d.getSubscribedSubdomainLabel(ctx, subaccountID)
	if err != nil {
		return err
	}

	regionLabel, err := d.getRegionLabel(ctx, *subdomainLabel.Tenant)
	if err != nil {
		return err
	}

	subdomain := subdomainLabel.Value.(string)
	region := regionLabel.Value.(string)

	instanceConfig, ok := d.destinationsConfig.RegionToInstanceConfig[region]
	if !ok {
		log.C(ctx).Errorf("No destination instance credentials found for region '%s'", region)
		return errors.New(fmt.Sprintf("No destination instance credentials found for region '%s'", region))
	}

	client, err := NewClient(instanceConfig, d.apiConfig, d.destinationsConfig.OauthTokenPath, subdomain)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to create Destination API client: %v", err)
		return err
	}

	if err := d.walkthroughPages(ctx, client, func(destinations []model.DestinationInput) error {
		log.C(ctx).Infof("Found %d destinations in subaccount '%s'", len(destinations), subaccountID)
		return d.mapDestinationsToTenant(ctx, *subdomainLabel.Tenant, destinations)
	}); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to sync destinations for subaccount '%s': %v", subaccountID, err)
		return err
	}

	return nil
}

func (d DestinationService) mapDestinationsToTenant(ctx context.Context, tenant string, destinations []model.DestinationInput) error {
	tx, err := d.transact.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to begin db transaction")
		return err
	}
	ctx = persistence.SaveToContext(ctx, tx)
	defer d.transact.RollbackUnlessCommitted(ctx, tx)

	for _, destination := range destinations {
		correlationID := correlationIDPrefix + destination.CommunicationScenarioId
		bundles, err := d.bundleRepo.GetBySystemAndCorrelationId(ctx, tenant, destination.XFSystemName, destination.URL, correlationID)

		if len(bundles) == 0 {
			log.C(ctx).Infof("No bundles found for system '%s', url '%s', correlation id '%s'", destination.XFSystemName, destination.URL, correlationID)
			continue
		}

		if err != nil {
			log.C(ctx).WithError(err).Errorf("Failed to fetch bundle for system '%s', url '%s', correlation id '%s', tenant id '%s': %v",
				destination.XFSystemName, destination.URL, correlationID, tenant, err)
			continue
		}

		for _, bundle := range bundles {
			id := d.uuidSvc.Generate()
			revision := d.uuidSvc.Generate()
			if err := d.repo.Upsert(ctx, destination, id, tenant, bundle.ID, revision); err != nil {
				log.C(ctx).WithError(err).Errorf("Failed to insert destination with name '%s' for bunlde '%s' and tenant '%s' to DB: %v", destination.Name, bundle.ID, tenant, err)
				continue
			}
		}
	}

	if err = tx.Commit(); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to commit database transaction %v", err)
		return err
	}
	return nil
}

type processFunc func([]model.DestinationInput) error

func (d DestinationService) walkthroughPages(ctx context.Context, client *Client, process processFunc) error {
	hasMorePages := true

	for page := 1; hasMorePages; page++ {
		pageString := strconv.Itoa(page)
		resp, err := client.FetchSubbacountDestinationsPage(ctx, pageString)
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

func (d DestinationService) FetchDestinationsSensitiveData(ctx context.Context, subaccountID string, destinationNames []string) ([]byte, error) {
	subdomainLabel, err := d.getSubscribedSubdomainLabel(ctx, subaccountID)
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

	instanceConfig, ok := d.destinationsConfig.RegionToInstanceConfig[region]
	if !ok {
		log.C(ctx).Errorf("No destination instance credentials found for region '%s'", region)
		return nil, errors.New(fmt.Sprintf("No destination instance credentials found for region '%s'", region))
	}
	client, err := NewClient(instanceConfig, d.apiConfig, d.destinationsConfig.OauthTokenPath, subdomain)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to create Destination API client: %v", err)
		return nil, err
	}

	nameCount := len(destinationNames)
	results := make([][]byte, nameCount)
	weighted := semaphore.NewWeighted(d.apiConfig.GoroutineLimit)
	resChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		for _, destination := range destinationNames {
			weighted.Acquire(ctx, 1)
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

	return append([]byte("{ \"destinations\": {"), combinedInfoJSON...), nil
}

func fetchDestination(ctx context.Context, dest string, weighted *semaphore.Weighted,
	client *Client, resChan chan []byte, errChan chan error) {

	log.C(ctx).Infof("Fetching data for destination: %s \n", dest)
	defer weighted.Release(1)
	result, err := client.FetchDestinationSensitiveData(ctx, dest)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to fetch data for destination %s: %v", dest, err)
		errChan <- err
		return
	}

	result = append([]byte("\""+dest+"\":"), result...)

	resChan <- result
}

func (d DestinationService) getSubscribedSubdomainLabel(ctx context.Context, subaccountID string) (*model.Label, error) {
	tx, err := d.transact.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to begin db transaction")
		return nil, err
	}
	ctx = persistence.SaveToContext(ctx, tx)
	defer d.transact.RollbackUnlessCommitted(ctx, tx)

	label, err := d.labelRepo.GetSubdomainLabelForSubscribedRuntime(ctx, subaccountID)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			log.C(ctx).Errorf("No subscribed subdomain found for subbaccount '%s'", subaccountID)
			return nil, apperrors.NewNotFoundErrorWithMessage(resource.Label, subaccountID, fmt.Sprintf("subaccount %s not found", subaccountID))
		}
		log.C(ctx).WithError(err).Errorf("Failed to get subdomain for subaccount '%s' from db: %v", subaccountID, err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to commit database transaction %v", err)
		return nil, err
	}

	return label, nil
}

func (d DestinationService) getRegionLabel(ctx context.Context, tenantID string) (*model.Label, error) {
	tx, err := d.transact.Begin()
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to begin db transaction")
		return nil, err
	}
	ctx = persistence.SaveToContext(ctx, tx)
	defer d.transact.RollbackUnlessCommitted(ctx, tx)

	region, err := d.labelRepo.GetByKey(ctx, tenantID, model.TenantLabelableObject, tenantID, regionLabelKey)
	if err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to fetch region for tenant '%s': %v", tenantID, err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.C(ctx).WithError(err).Errorf("Failed to commit database transaction %v", err)
		return nil, err
	}

	return region, nil
}
