package bundle

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kyma-incubator/compass/components/director/internal/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal/labelfilter"

	"github.com/kyma-incubator/compass/components/director/pkg/pagination"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
	"github.com/pkg/errors"
)

const bundleTable string = `public.bundles`

var (
	bundleColumns    = []string{"id", "tenant_id", "app_id", "name", "description", "instance_auth_request_json_schema", "default_instance_auth", "ord_id", "short_description", "links", "labels", "credential_exchange_strategies", "ready", "created_at", "updated_at", "deleted_at", "error"}
	tenantColumn     = "tenant_id"
	updatableColumns = []string{"name", "description", "instance_auth_request_json_schema", "default_instance_auth", "ord_id", "short_description", "links", "labels", "credential_exchange_strategies", "ready", "created_at", "updated_at", "deleted_at", "error"}
	orderByColumns   = repo.OrderByParams{repo.NewAscOrderBy("app_id"), repo.NewAscOrderBy("id")}
)

// EntityConverter missing godoc
//go:generate mockery --name=EntityConverter --output=automock --outpkg=automock --case=underscore
type EntityConverter interface {
	ToEntity(in *model.Bundle) (*Entity, error)
	FromEntity(entity *Entity) (*model.Bundle, error)
}

type pgRepository struct {
	existQuerier repo.ExistQuerier
	singleGetter repo.SingleGetter
	deleter      repo.Deleter
	lister       repo.Lister
	unionLister  repo.UnionLister
	creator      repo.Creator
	updater      repo.Updater
	conv         EntityConverter
}

// NewRepository missing godoc
func NewRepository(conv EntityConverter) *pgRepository {
	return &pgRepository{
		existQuerier: repo.NewExistQuerier(resource.Bundle, bundleTable, tenantColumn),
		singleGetter: repo.NewSingleGetter(resource.Bundle, bundleTable, tenantColumn, bundleColumns),
		deleter:      repo.NewDeleter(resource.Bundle, bundleTable, tenantColumn),
		lister:       repo.NewLister(resource.Bundle, bundleTable, tenantColumn, bundleColumns),
		unionLister:  repo.NewUnionLister(resource.Bundle, bundleTable, tenantColumn, bundleColumns),
		creator:      repo.NewCreator(resource.Bundle, bundleTable, bundleColumns),
		updater:      repo.NewUpdater(resource.Bundle, bundleTable, updatableColumns, tenantColumn, []string{"id"}),
		conv:         conv,
	}
}

// BundleCollection missing godoc
type BundleCollection []Entity

// Len missing godoc
func (r BundleCollection) Len() int {
	return len(r)
}

// Create missing godoc
func (r *pgRepository) Create(ctx context.Context, model *model.Bundle) error {
	if model == nil {
		return apperrors.NewInternalError("model can not be nil")
	}

	bndlEnt, err := r.conv.ToEntity(model)
	if err != nil {
		return errors.Wrap(err, "while converting to Bundle entity")
	}

	log.C(ctx).Debugf("Persisting Bundle entity with id %s to db", model.ID)
	return r.creator.Create(ctx, bndlEnt)
}

// Update missing godoc
func (r *pgRepository) Update(ctx context.Context, model *model.Bundle) error {
	if model == nil {
		return apperrors.NewInternalError("model can not be nil")
	}

	bndlEnt, err := r.conv.ToEntity(model)

	if err != nil {
		return errors.Wrap(err, "while converting to Bundle entity")
	}

	return r.updater.UpdateSingle(ctx, bndlEnt)
}

// Delete missing godoc
func (r *pgRepository) Delete(ctx context.Context, tenant, id string) error {
	return r.deleter.DeleteOne(ctx, tenant, repo.Conditions{repo.NewEqualCondition("id", id)})
}

// Exists missing godoc
func (r *pgRepository) Exists(ctx context.Context, tenant, id string) (bool, error) {
	return r.existQuerier.Exists(ctx, tenant, repo.Conditions{repo.NewEqualCondition("id", id)})
}

// GetByID missing godoc
func (r *pgRepository) GetByID(ctx context.Context, tenant, id string) (*model.Bundle, error) {
	var bndlEnt Entity
	if err := r.singleGetter.Get(ctx, tenant, repo.Conditions{repo.NewEqualCondition("id", id)}, repo.NoOrderBy, &bndlEnt); err != nil {
		return nil, err
	}

	bndlModel, err := r.conv.FromEntity(&bndlEnt)
	if err != nil {
		return nil, errors.Wrap(err, "while converting Bundle from Entity")
	}

	return bndlModel, nil
}

// GetForApplication missing godoc
func (r *pgRepository) GetForApplication(ctx context.Context, tenant string, id string, applicationID string) (*model.Bundle, error) {
	var ent Entity

	conditions := repo.Conditions{
		repo.NewEqualCondition("id", id),
		repo.NewEqualCondition("app_id", applicationID),
	}
	if err := r.singleGetter.Get(ctx, tenant, conditions, repo.NoOrderBy, &ent); err != nil {
		return nil, err
	}

	bndlModel, err := r.conv.FromEntity(&ent)
	if err != nil {
		return nil, errors.Wrap(err, "while creating Bundle model from entity")
	}

	return bndlModel, nil
}

// ListByApplicationIDs retrieves a page of bundles for a given number of application IDs.
func (r *pgRepository) ListByApplicationIDs(ctx context.Context, tenantID string, applicationIDs []string, pageSize int, cursor string) ([]*model.BundlePage, error) {
	var bundleCollection BundleCollection
	counts, err := r.unionLister.List(ctx, tenantID, applicationIDs, "app_id", pageSize, cursor, orderByColumns, &bundleCollection)
	if err != nil {
		return nil, err
	}

	bundleByID := map[string][]*model.Bundle{}
	for _, bundleEnt := range bundleCollection {
		m, err := r.conv.FromEntity(&bundleEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating Bundle model from entity")
		}
		bundleByID[bundleEnt.ApplicationID] = append(bundleByID[bundleEnt.ApplicationID], m)
	}

	offset, err := pagination.DecodeOffsetCursor(cursor)
	if err != nil {
		return nil, errors.Wrap(err, "while decoding page cursor")
	}

	bundlePages := make([]*model.BundlePage, 0, len(applicationIDs))
	for _, appID := range applicationIDs {
		totalCount := counts[appID]
		hasNextPage := false
		endCursor := ""
		if totalCount > offset+len(bundleByID[appID]) {
			hasNextPage = true
			endCursor = pagination.EncodeNextOffsetCursor(offset, pageSize)
		}

		page := &pagination.Page{
			StartCursor: cursor,
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		}

		bundlePages = append(bundlePages, &model.BundlePage{Data: bundleByID[appID], TotalCount: totalCount, PageInfo: page})
	}

	return bundlePages, nil
}

// ListByApplicationIDNoPaging missing godoc
func (r *pgRepository) ListByApplicationIDNoPaging(ctx context.Context, tenantID, appID string) ([]*model.Bundle, error) {
	bundleCollection := BundleCollection{}
	if err := r.lister.List(ctx, tenantID, &bundleCollection, repo.NewEqualCondition("app_id", appID)); err != nil {
		return nil, err
	}
	bundles := make([]*model.Bundle, 0, bundleCollection.Len())
	for _, bundle := range bundleCollection {
		bundleModel, err := r.conv.FromEntity(&bundle)
		if err != nil {
			return nil, err
		}
		bundles = append(bundles, bundleModel)
	}
	return bundles, nil
}

// ListByApplicationIDsForScenarios retrieves a page of bundles that are part of at least one of the provided scenarios for a given number of application IDs.
func (r *pgRepository) ListByApplicationIDsForScenarios(ctx context.Context, tenantID string, applicationIDs []string, scenarios []string, pageSize int, cursor string) ([]*model.BundlePage, error) {
	var bundleCollection BundleCollection

	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, apperrors.NewInvalidDataError("tenantID is not UUID")
	}

	scenariosFilters := make([]*labelfilter.LabelFilter, 0, len(scenarios))
	for _, scenarioValue := range scenarios {
		query := fmt.Sprintf(`$[*] ? (@ == "%s")`, scenarioValue)
		scenariosFilters = append(scenariosFilters, labelfilter.NewForKeyWithQuery(model.ScenariosKey, query))
	}

	scenariosSubquery, scenariosArgs, err := label.FilterQuery(model.BundleLabelableObject, label.UnionSet, tenantUUID, scenariosFilters)
	if err != nil {
		return nil, errors.Wrap(err, "while creating scenarios filter query")
	}

	counts, err := r.unionLister.List(ctx,
		tenantID,
		applicationIDs,
		"app_id",
		pageSize,
		cursor,
		orderByColumns,
		&bundleCollection,
		repo.NewInConditionForSubQuery("id", scenariosSubquery, scenariosArgs))
	if err != nil {
		return nil, err
	}

	bundleByID := map[string][]*model.Bundle{}
	for _, bundleEnt := range bundleCollection {
		m, err := r.conv.FromEntity(&bundleEnt)
		if err != nil {
			return nil, errors.Wrap(err, "while creating Bundle model from entity")
		}
		bundleByID[bundleEnt.ApplicationID] = append(bundleByID[bundleEnt.ApplicationID], m)
	}

	offset, err := pagination.DecodeOffsetCursor(cursor)
	if err != nil {
		return nil, errors.Wrap(err, "while decoding page cursor")
	}

	bundlePages := make([]*model.BundlePage, 0, len(applicationIDs))
	for _, appID := range applicationIDs {
		totalCount := counts[appID]
		hasNextPage := false
		endCursor := ""
		if totalCount > offset+len(bundleByID[appID]) {
			hasNextPage = true
			endCursor = pagination.EncodeNextOffsetCursor(offset, pageSize)
		}

		page := &pagination.Page{
			StartCursor: cursor,
			EndCursor:   endCursor,
			HasNextPage: hasNextPage,
		}

		bundlePages = append(bundlePages, &model.BundlePage{Data: bundleByID[appID], TotalCount: totalCount, PageInfo: page})
	}

	return bundlePages, nil
}

// ListByApplicationIDsForScenariosNoPaging retrieves bundles part of a set of applications that are part of at least one of the provided scenarios.
func (r *pgRepository) ListByApplicationIDsForScenariosNoPaging(ctx context.Context, tenantID string, applicationIDs []string, scenarios []string) ([]*model.Bundle, error) {
	var bundleCollection BundleCollection

	tenantUUID, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, apperrors.NewInvalidDataError("tenantID is not UUID")
	}

	scenariosFilters := make([]*labelfilter.LabelFilter, 0, len(scenarios))
	for _, scenarioValue := range scenarios {
		query := fmt.Sprintf(`$[*] ? (@ == "%s")`, scenarioValue)
		scenariosFilters = append(scenariosFilters, labelfilter.NewForKeyWithQuery(model.ScenariosKey, query))
	}

	scenariosSubquery, scenariosArgs, err := label.FilterQuery(model.BundleLabelableObject, label.UnionSet, tenantUUID, scenariosFilters)
	if err != nil {
		return nil, errors.Wrap(err, "while creating scenarios filter query")
	}

	if err := r.lister.List(ctx,
		tenantID,
		&bundleCollection,
		repo.NewInConditionForStringValues("app_id", applicationIDs),
		repo.NewInConditionForSubQuery("id", scenariosSubquery, scenariosArgs)); err != nil {
		return nil, err
	}

	bundles := make([]*model.Bundle, 0, bundleCollection.Len())
	for _, bundle := range bundleCollection {
		bundleModel, err := r.conv.FromEntity(&bundle)
		if err != nil {
			return nil, err
		}
		bundles = append(bundles, bundleModel)
	}
	return bundles, nil
}
