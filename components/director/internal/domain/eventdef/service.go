package eventdef

import (
	"context"

	"github.com/kyma-incubator/compass/components/director/pkg/str"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/timestamp"

	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/pkg/errors"
)

// EventAPIRepository is responsible for the repo-layer EventDefinition operations.
//
//go:generate mockery --name=EventAPIRepository --output=automock --outpkg=automock --case=underscore --disable-version-string
type EventAPIRepository interface {
	GetByID(ctx context.Context, tenantID string, id string) (*model.EventDefinition, error)
	GetByIDGlobal(ctx context.Context, id string) (*model.EventDefinition, error)
	GetForBundle(ctx context.Context, tenant string, id string, bundleID string) (*model.EventDefinition, error)
	GetByApplicationID(ctx context.Context, tenantID string, id, appID string) (*model.EventDefinition, error)
	ListByBundleIDs(ctx context.Context, tenantID string, bundleIDs []string, bundleRefs []*model.BundleReference, totalCounts map[string]int, pageSize int, cursor string) ([]*model.EventDefinitionPage, error)
	ListByResourceID(ctx context.Context, tenantID, resourceID string, resourceType resource.Type) ([]*model.EventDefinition, error)
	ListByApplicationIDPage(ctx context.Context, tenantID string, appID string, pageSize int, cursor string) (*model.EventDefinitionPage, error)
	Create(ctx context.Context, tenant string, item *model.EventDefinition) error
	CreateGlobal(ctx context.Context, item *model.EventDefinition) error
	Update(ctx context.Context, tenant string, item *model.EventDefinition) error
	UpdateGlobal(ctx context.Context, item *model.EventDefinition) error
	Delete(ctx context.Context, tenantID string, id string) error
	DeleteGlobal(ctx context.Context, id string) error
	DeleteAllByBundleID(ctx context.Context, tenantID, bundleID string) error
}

// UIDService is responsible for generating GUIDs, which will be used as internal eventDefinition IDs when they are created.
//
//go:generate mockery --name=UIDService --output=automock --outpkg=automock --case=underscore --disable-version-string
type UIDService interface {
	Generate() string
}

// SpecService is responsible for the service-layer Specification operations.
//
//go:generate mockery --name=SpecService --output=automock --outpkg=automock --case=underscore --disable-version-string
type SpecService interface {
	CreateByReferenceObjectID(ctx context.Context, in model.SpecInput, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) (string, error)
	UpdateByReferenceObjectID(ctx context.Context, id string, in model.SpecInput, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) error
	GetByReferenceObjectID(ctx context.Context, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) (*model.Spec, error)
	RefetchSpec(ctx context.Context, id string, objectType model.SpecReferenceObjectType) (*model.Spec, error)
	ListFetchRequestsByReferenceObjectIDs(ctx context.Context, tenant string, objectIDs []string, objectType model.SpecReferenceObjectType) ([]*model.FetchRequest, error)
}

// BundleReferenceService is responsible for the service-layer BundleReference operations.
//
//go:generate mockery --name=BundleReferenceService --output=automock --outpkg=automock --case=underscore --disable-version-string
type BundleReferenceService interface {
	GetForBundle(ctx context.Context, objectType model.BundleReferenceObjectType, objectID, bundleID *string) (*model.BundleReference, error)
	CreateByReferenceObjectID(ctx context.Context, in model.BundleReferenceInput, objectType model.BundleReferenceObjectType, objectID, bundleID *string) error
	UpdateByReferenceObjectID(ctx context.Context, in model.BundleReferenceInput, objectType model.BundleReferenceObjectType, objectID, bundleID *string) error
	DeleteByReferenceObjectID(ctx context.Context, objectType model.BundleReferenceObjectType, objectID, bundleID *string) error
	ListByBundleIDs(ctx context.Context, objectType model.BundleReferenceObjectType, bundleIDs []string, pageSize int, cursor string) ([]*model.BundleReference, map[string]int, error)
}

type service struct {
	eventAPIRepo           EventAPIRepository
	uidService             UIDService
	specService            SpecService
	bundleReferenceService BundleReferenceService
	timestampGen           timestamp.Generator
}

// NewService returns a new object responsible for service-layer EventDefinition operations.
func NewService(eventAPIRepo EventAPIRepository, uidService UIDService, specService SpecService, bundleReferenceService BundleReferenceService) *service {
	return &service{
		eventAPIRepo:           eventAPIRepo,
		uidService:             uidService,
		specService:            specService,
		bundleReferenceService: bundleReferenceService,
		timestampGen:           timestamp.DefaultGenerator,
	}
}

// ListByBundleIDs lists all EventDefinitions in pages for a given array of bundle IDs.
func (s *service) ListByBundleIDs(ctx context.Context, bundleIDs []string, pageSize int, cursor string) ([]*model.EventDefinitionPage, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "while loading tenant from context")
	}

	if pageSize < 1 || pageSize > 200 {
		return nil, apperrors.NewInvalidDataError("page size must be between 1 and 200")
	}

	bundleRefs, counts, err := s.bundleReferenceService.ListByBundleIDs(ctx, model.BundleEventReference, bundleIDs, pageSize, cursor)
	if err != nil {
		return nil, err
	}

	return s.eventAPIRepo.ListByBundleIDs(ctx, tnt, bundleIDs, bundleRefs, counts, pageSize, cursor)
}

// ListByApplicationID lists all EventDefinitions for a given application ID.
func (s *service) ListByApplicationID(ctx context.Context, appID string) ([]*model.EventDefinition, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.eventAPIRepo.ListByResourceID(ctx, tnt, appID, resource.Application)
}

// ListByApplicationTemplateVersionID lists all EventDefinitions for a given application template version ID.
func (s *service) ListByApplicationTemplateVersionID(ctx context.Context, appTemplateVersionID string) ([]*model.EventDefinition, error) {
	return s.eventAPIRepo.ListByResourceID(ctx, "", appTemplateVersionID, resource.ApplicationTemplateVersion)
}

// ListByApplicationIDPage lists all EventDefinitions for a given application ID with paging.
func (s *service) ListByApplicationIDPage(ctx context.Context, appID string, pageSize int, cursor string) (*model.EventDefinitionPage, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if pageSize < 1 || pageSize > 200 {
		return nil, apperrors.NewInvalidDataError("page size must be between 1 and 200")
	}

	return s.eventAPIRepo.ListByApplicationIDPage(ctx, tnt, appID, pageSize, cursor)
}

// Get returns the EventDefinition by its ID.
func (s *service) Get(ctx context.Context, id string) (*model.EventDefinition, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	eventAPI, err := s.eventAPIRepo.GetByID(ctx, tnt, id)
	if err != nil {
		return nil, err
	}

	return eventAPI, nil
}

// GetForBundle returns an EventDefinition by its ID and a bundle ID.
func (s *service) GetForBundle(ctx context.Context, id string, bundleID string) (*model.EventDefinition, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	eventAPI, err := s.eventAPIRepo.GetForBundle(ctx, tnt, id, bundleID)
	if err != nil {
		return nil, errors.Wrap(err, "while getting API definition")
	}

	return eventAPI, nil
}

// GetForApplication returns an EventDefinition by its ID and Application ID.
func (s *service) GetForApplication(ctx context.Context, id string, appID string) (*model.EventDefinition, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	eventDefinition, err := s.eventAPIRepo.GetByApplicationID(ctx, tnt, id, appID)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting Event Definition with id %q", id)
	}

	return eventDefinition, nil
}

// CreateInBundle creates an EventDefinition. This function is used in the graphQL flow.
func (s *service) CreateInBundle(ctx context.Context, resourceType resource.Type, resourceID string, bundleID string, in model.EventDefinitionInput, spec *model.SpecInput) (string, error) {
	return s.Create(ctx, resourceType, resourceID, &bundleID, nil, in, []*model.SpecInput{spec}, nil, 0, "")
}

// CreateInApplication creates an EventDefinition in the context of an Application without Bundle
func (s *service) CreateInApplication(ctx context.Context, appID string, in model.EventDefinitionInput, spec *model.SpecInput) (string, error) {
	return s.Create(ctx, resource.Application, appID, nil, nil, in, []*model.SpecInput{spec}, nil, 0, "")
}

// Create creates EventDefinition/s. This function is used both in the ORD scenario and is re-used in CreateInBundle but with "null" ORD specific arguments.
func (s *service) Create(ctx context.Context, resourceType resource.Type, resourceID string, bundleID, packageID *string, in model.EventDefinitionInput, specs []*model.SpecInput, bundleIDs []string, eventHash uint64, defaultBundleID string) (string, error) {
	id := s.uidService.Generate()
	eventAPI := in.ToEventDefinition(id, resourceType, resourceID, packageID, eventHash)

	if err := s.createEventDef(ctx, eventAPI, resourceType); err != nil {
		return "", errors.Wrap(err, "while creating api")
	}

	if err := s.processSpecs(ctx, eventAPI.ID, specs, resourceType); err != nil {
		return "", err
	}

	if err := s.createBundleReferenceObject(ctx, eventAPI.ID, bundleID, defaultBundleID, bundleIDs); err != nil {
		return "", err
	}

	return id, nil
}

// Update updates an EventDefinition. This function is used in the graphQL flow.
func (s *service) Update(ctx context.Context, resourceType resource.Type, id string, in model.EventDefinitionInput, specIn *model.SpecInput) error {
	return s.UpdateInManyBundles(ctx, resourceType, id, in, specIn, nil, nil, nil, 0, "")
}

// UpdateInManyBundles updates EventDefinition/s. This function is used both in the ORD scenario and is re-used in Update but with "null" ORD specific arguments.
func (s *service) UpdateInManyBundles(ctx context.Context, resourceType resource.Type, id string, in model.EventDefinitionInput, specIn *model.SpecInput, bundleIDsFromBundleReference, bundleIDsForCreation, bundleIDsForDeletion []string, eventHash uint64, defaultBundleID string) error {
	eventDef, err := s.getEventDef(ctx, id, resourceType)
	if err != nil {
		return errors.Wrapf(err, "while getting EventDefinition with ID %s", id)
	}

	resourceID := getParentResourceID(eventDef)
	eventDef = in.ToEventDefinition(id, resourceType, resourceID, eventDef.PackageID, eventHash)

	if err = s.updateEventDef(ctx, eventDef, resourceType); err != nil {
		return errors.Wrapf(err, "while updating EventDefinition with ID %s", id)
	}

	if err = s.handleReferenceObjects(ctx, eventDef.ID, bundleIDsForCreation, bundleIDsForDeletion, bundleIDsFromBundleReference, defaultBundleID); err != nil {
		return errors.Wrapf(err, "while handling reference objects for EventDefinition with ID %s", eventDef.ID)
	}

	if specIn != nil {
		return s.handleSpecsInEvent(ctx, resourceType, eventDef.ID, specIn)
	}

	return nil
}

// UpdateForApplication updates an EventDefinition for Application without being in a Bundle
func (s *service) UpdateForApplication(ctx context.Context, id string, in model.EventDefinitionInput, specIn *model.SpecInput) error {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	event, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	event = in.ToEventDefinition(id, resource.Application, str.PtrStrToStr(event.ApplicationID), event.PackageID, 0)

	if err = s.eventAPIRepo.Update(ctx, tnt, event); err != nil {
		return errors.Wrapf(err, "while updating EventDefinition with id %s", id)
	}

	if specIn != nil {
		return s.handleSpecsInEvent(ctx, resource.Application, id, specIn)
	}

	return nil
}

// Delete deletes the EventDefinition by its ID.
func (s *service) Delete(ctx context.Context, resourceType resource.Type, id string) error {
	if err := s.deleteEventDef(ctx, id, resourceType); err != nil {
		return errors.Wrapf(err, "while deleting EventDefinition with id %s", id)
	}

	log.C(ctx).Infof("Successfully deleted EventDefinition with id %s", id)

	return nil
}

// DeleteAllByBundleID deletes all EventDefinitions for a given bundle ID
func (s *service) DeleteAllByBundleID(ctx context.Context, bundleID string) error {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	err = s.eventAPIRepo.DeleteAllByBundleID(ctx, tnt, bundleID)
	if err != nil {
		return errors.Wrapf(err, "while deleting EventDefinitions for Bundle with id %q", bundleID)
	}

	return nil
}

// ListFetchRequests lists all FetchRequests for given specification IDs
func (s *service) ListFetchRequests(ctx context.Context, specIDs []string) ([]*model.FetchRequest, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	fetchRequests, err := s.specService.ListFetchRequestsByReferenceObjectIDs(ctx, tnt, specIDs, model.EventSpecReference)
	if err != nil {
		if apperrors.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return fetchRequests, nil
}

func (s *service) handleSpecsInEvent(ctx context.Context, resourceType resource.Type, id string, specIn *model.SpecInput) error {
	dbSpec, err := s.specService.GetByReferenceObjectID(ctx, resourceType, model.EventSpecReference, id)
	if err != nil {
		return errors.Wrapf(err, "while getting spec for EventDefinition with id %q", id)
	}

	if dbSpec == nil {
		_, err = s.specService.CreateByReferenceObjectID(ctx, *specIn, resourceType, model.EventSpecReference, id)
		return errors.Wrapf(err, "while creating spec for EventDefinition with ID %s", id)
	}

	if err = s.specService.UpdateByReferenceObjectID(ctx, dbSpec.ID, *specIn, resourceType, model.EventSpecReference, id); err != nil {
		return errors.Wrapf(err, "while updating spec for EventDefinition with ID %s", id)
	}

	return nil
}

func (s *service) handleReferenceObjects(ctx context.Context, id string, bundleIDsForCreation, bundleIDsForDeletion, bundleIDsFromBundleReference []string, defaultBundleID string) error {
	for _, bundleID := range bundleIDsForCreation {
		createBundleRefInput := &model.BundleReferenceInput{}
		if defaultBundleID != "" && bundleID == defaultBundleID {
			isDefaultBundle := true
			createBundleRefInput = &model.BundleReferenceInput{IsDefaultBundle: &isDefaultBundle}
		}
		if err := s.bundleReferenceService.CreateByReferenceObjectID(ctx, *createBundleRefInput, model.BundleEventReference, &id, &bundleID); err != nil {
			return err
		}
	}

	for _, bundleID := range bundleIDsForDeletion {
		if err := s.bundleReferenceService.DeleteByReferenceObjectID(ctx, model.BundleEventReference, &id, &bundleID); err != nil {
			return err
		}
	}

	for _, bundleID := range bundleIDsFromBundleReference {
		bundleRefInput := &model.BundleReferenceInput{}
		if defaultBundleID != "" && bundleID == defaultBundleID {
			isDefaultBundle := true
			bundleRefInput = &model.BundleReferenceInput{IsDefaultBundle: &isDefaultBundle}
		}
		if err := s.bundleReferenceService.UpdateByReferenceObjectID(ctx, *bundleRefInput, model.BundleEventReference, &id, &bundleID); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) processSpecs(ctx context.Context, eventID string, specs []*model.SpecInput, resourceType resource.Type) error {
	for _, spec := range specs {
		if spec == nil {
			continue
		}

		if _, err := s.specService.CreateByReferenceObjectID(ctx, *spec, resourceType, model.EventSpecReference, eventID); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) createBundleReferenceObject(ctx context.Context, eventID string, bundleID *string, defaultBundleID string, bundleIDs []string) error {
	if bundleIDs == nil && bundleID != nil {
		if err := s.bundleReferenceService.CreateByReferenceObjectID(ctx, model.BundleReferenceInput{}, model.BundleEventReference, &eventID, bundleID); err != nil {
			return err
		}
	}

	for _, bndlID := range bundleIDs {
		bundleRefInput := &model.BundleReferenceInput{}
		if defaultBundleID != "" && bndlID == defaultBundleID {
			isDefaultBundle := true
			bundleRefInput = &model.BundleReferenceInput{
				IsDefaultBundle: &isDefaultBundle,
			}
		}
		if err := s.bundleReferenceService.CreateByReferenceObjectID(ctx, *bundleRefInput, model.BundleEventReference, &eventID, &bndlID); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) createEventDef(ctx context.Context, eventAPI *model.EventDefinition, resourceType resource.Type) error {
	if resourceType.IsTenantIgnorable() {
		return s.eventAPIRepo.CreateGlobal(ctx, eventAPI)
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	return s.eventAPIRepo.Create(ctx, tnt, eventAPI)
}

func (s *service) getEventDef(ctx context.Context, id string, resourceType resource.Type) (*model.EventDefinition, error) {
	if resourceType.IsTenantIgnorable() {
		return s.eventAPIRepo.GetByIDGlobal(ctx, id)
	}

	return s.Get(ctx, id)
}

func (s *service) updateEventDef(ctx context.Context, eventDef *model.EventDefinition, resourceType resource.Type) error {
	if resourceType.IsTenantIgnorable() {
		return s.eventAPIRepo.UpdateGlobal(ctx, eventDef)
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	return s.eventAPIRepo.Update(ctx, tnt, eventDef)
}

func (s *service) deleteEventDef(ctx context.Context, id string, resourceType resource.Type) error {
	if resourceType.IsTenantIgnorable() {
		return s.eventAPIRepo.DeleteGlobal(ctx, id)
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return errors.Wrapf(err, "while loading tenant from context")
	}

	return s.eventAPIRepo.Delete(ctx, tnt, id)
}

func getParentResourceID(api *model.EventDefinition) string {
	if api.ApplicationTemplateVersionID != nil {
		return *api.ApplicationTemplateVersionID
	} else if api.ApplicationID != nil {
		return *api.ApplicationID
	}

	return ""
}
