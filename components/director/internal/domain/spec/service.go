package spec

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/pkg/resource"

	"github.com/kyma-incubator/compass/components/director/internal/domain/tenant"
	"github.com/kyma-incubator/compass/components/director/internal/timestamp"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/pkg/errors"

	"github.com/kyma-incubator/compass/components/director/internal/model"
)

// SpecRepository missing godoc
//
//go:generate mockery --name=SpecRepository --output=automock --outpkg=automock --case=underscore --disable-version-string
type SpecRepository interface {
	Create(ctx context.Context, tenant string, item *model.Spec) error
	CreateGlobal(ctx context.Context, item *model.Spec) error
	GetByID(ctx context.Context, tenantID string, id string, objectType model.SpecReferenceObjectType) (*model.Spec, error)
	GetByIDGlobal(ctx context.Context, id string) (*model.Spec, error)
	ListIDByReferenceObjectID(ctx context.Context, tenant string, objectType model.SpecReferenceObjectType, objectID string) ([]string, error)
	ListIDByReferenceObjectIDGlobal(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string) ([]string, error)
	ListByReferenceObjectID(ctx context.Context, tenant string, objectType model.SpecReferenceObjectType, objectID string) ([]*model.Spec, error)
	ListByReferenceObjectIDGlobal(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string) ([]*model.Spec, error)
	ListByReferenceObjectIDs(ctx context.Context, tenant string, objectType model.SpecReferenceObjectType, objectIDs []string) ([]*model.Spec, error)
	Delete(ctx context.Context, tenant, id string, objectType model.SpecReferenceObjectType) error
	DeleteByReferenceObjectID(ctx context.Context, tenant string, objectType model.SpecReferenceObjectType, objectID string) error
	DeleteByReferenceObjectIDGlobal(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string) error
	Update(ctx context.Context, tenant string, item *model.Spec) error
	UpdateGlobal(ctx context.Context, item *model.Spec) error
	Exists(ctx context.Context, tenantID, id string, objectType model.SpecReferenceObjectType) (bool, error)
}

// FetchRequestRepository missing godoc
//
//go:generate mockery --name=FetchRequestRepository --output=automock --outpkg=automock --case=underscore --disable-version-string
type FetchRequestRepository interface {
	Create(ctx context.Context, tenant string, item *model.FetchRequest) error
	CreateGlobal(ctx context.Context, item *model.FetchRequest) error
	GetByReferenceObjectID(ctx context.Context, tenant string, objectType model.FetchRequestReferenceObjectType, objectID string) (*model.FetchRequest, error)
	DeleteByReferenceObjectID(ctx context.Context, tenant string, objectType model.FetchRequestReferenceObjectType, objectID string) error
	DeleteByReferenceObjectIDGlobal(ctx context.Context, objectType model.FetchRequestReferenceObjectType, objectID string) error
	ListByReferenceObjectIDs(ctx context.Context, tenant string, objectType model.FetchRequestReferenceObjectType, objectIDs []string) ([]*model.FetchRequest, error)
	ListByReferenceObjectIDsGlobal(ctx context.Context, objectType model.FetchRequestReferenceObjectType, objectIDs []string) ([]*model.FetchRequest, error)
}

// UIDService missing godoc
//
//go:generate mockery --name=UIDService --output=automock --outpkg=automock --case=underscore --disable-version-string
type UIDService interface {
	Generate() string
}

// FetchRequestService missing godoc
//
//go:generate mockery --name=FetchRequestService --output=automock --outpkg=automock --case=underscore --disable-version-string
type FetchRequestService interface {
	HandleSpec(ctx context.Context, fr *model.FetchRequest) *string
}

type service struct {
	repo                SpecRepository
	fetchRequestRepo    FetchRequestRepository
	uidService          UIDService
	fetchRequestService FetchRequestService
	timestampGen        timestamp.Generator
}

// NewService missing godoc
func NewService(repo SpecRepository, fetchRequestRepo FetchRequestRepository, uidService UIDService, fetchRequestService FetchRequestService) *service {
	return &service{
		repo:                repo,
		fetchRequestRepo:    fetchRequestRepo,
		uidService:          uidService,
		fetchRequestService: fetchRequestService,
		timestampGen:        timestamp.DefaultGenerator,
	}
}

// GetByID takes care of retrieving a specific spec entity from db based on a provided id and objectType (API or Event)
func (s *service) GetByID(ctx context.Context, id string, objectType model.SpecReferenceObjectType) (*model.Spec, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, tnt, id, objectType)
}

func (s *service) GetByIDGlobal(ctx context.Context, id string) (*model.Spec, error) {
	return s.repo.GetByIDGlobal(ctx, id)
}

// ListByReferenceObjectID missing godoc
func (s *service) ListByReferenceObjectID(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string) ([]*model.Spec, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.ListByReferenceObjectID(ctx, tnt, objectType, objectID)
}

// ListIDByReferenceObjectID retrieves all spec ids by objectType and objectID
func (s *service) ListIDByReferenceObjectID(ctx context.Context, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) ([]string, error) {
	if resourceType.IsTenantIgnorable() {
		return s.repo.ListIDByReferenceObjectIDGlobal(ctx, objectType, objectID)
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.ListIDByReferenceObjectID(ctx, tnt, objectType, objectID)
}

// GetByReferenceObjectID
// Until now APIs and Events had embedded specification in them, we will model this behavior by relying that the first created spec is the one which GraphQL expects
func (s *service) GetByReferenceObjectID(ctx context.Context, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) (*model.Spec, error) {
	specs, err := s.listSpecsByReferenceObjectID(ctx, objectType, objectID, resourceType)
	if err != nil {
		return nil, err
	}

	if len(specs) > 0 {
		return specs[0], nil
	}

	return nil, nil
}

// ListByReferenceObjectIDs missing godoc
func (s *service) ListByReferenceObjectIDs(ctx context.Context, objectType model.SpecReferenceObjectType, objectIDs []string) ([]*model.Spec, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	specs, err := s.repo.ListByReferenceObjectIDs(ctx, tnt, objectType, objectIDs)

	return specs, err
}

// CreateByReferenceObjectID missing godoc
func (s *service) CreateByReferenceObjectID(ctx context.Context, in model.SpecInput, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) (string, error) {
	id := s.uidService.Generate()
	spec, err := in.ToSpec(id, objectType, objectID)
	if err != nil {
		return "", err
	}

	if err = s.createSpec(ctx, spec, resourceType); err != nil {
		return "", errors.Wrapf(err, "while creating spec for %q with id %q", objectType, objectID)
	}

	if in.Data == nil && in.FetchRequest != nil {
		fr, err := s.createFetchRequest(ctx, *in.FetchRequest, id, objectType, resourceType)
		if err != nil {
			return "", errors.Wrapf(err, "while creating FetchRequest for %s Specification with id %q", objectType, id)
		}

		spec.Data = s.fetchRequestService.HandleSpec(ctx, fr)

		if err = s.updateSpec(ctx, spec, resourceType); err != nil {
			return "", errors.Wrapf(err, "while updating %s Specification with id %q", objectType, id)
		}
	}

	return id, nil
}

// CreateByReferenceObjectIDWithDelayedFetchRequest identical to CreateByReferenceObjectID with the only difference that the spec and fetch request entities are only persisted in DB and the fetch request itself is not executed
func (s *service) CreateByReferenceObjectIDWithDelayedFetchRequest(ctx context.Context, in model.SpecInput, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) (string, *model.FetchRequest, error) {
	id := s.uidService.Generate()
	spec, err := in.ToSpec(id, objectType, objectID)
	if err != nil {
		return "", nil, err
	}

	if err = s.createSpec(ctx, spec, resourceType); err != nil {
		return "", nil, errors.Wrapf(err, "while creating spec for %q with id %q", objectType, objectID)
	}

	var fr *model.FetchRequest
	if in.Data == nil && in.FetchRequest != nil {
		fr, err = s.createFetchRequest(ctx, *in.FetchRequest, id, objectType, resourceType)
		if err != nil {
			return "", nil, errors.Wrapf(err, "while creating FetchRequest for %s Specification with id %q", objectType, id)
		}
	}

	return id, fr, nil
}

// UpdateByReferenceObjectID missing godoc
func (s *service) UpdateByReferenceObjectID(ctx context.Context, id string, in model.SpecInput, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) error {
	if err := s.deleteFetchRequestByReferenceObjectID(ctx, id, objectType, resourceType); err != nil {
		return errors.Wrapf(err, "while deleting FetchRequest for Specification with id %q", id)
	}

	spec, err := in.ToSpec(id, objectType, objectID)
	if err != nil {
		return err
	}

	if in.Data == nil && in.FetchRequest != nil {
		fr, err := s.createFetchRequest(ctx, *in.FetchRequest, id, objectType, resourceType)
		if err != nil {
			return errors.Wrapf(err, "while creating FetchRequest for %s Specification with id %q", objectType, id)
		}

		spec.Data = s.fetchRequestService.HandleSpec(ctx, fr)
	}

	if err = s.updateSpec(ctx, spec, resourceType); err != nil {
		return errors.Wrapf(err, "while updating %s Specification with id %q", objectType, id)
	}

	return nil
}

// UpdateSpecOnly takes care of simply updating a single spec entity in db without looking and executing corresponding fetch requests that may be related to it
func (s *service) UpdateSpecOnly(ctx context.Context, spec model.Spec) error {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	if err = s.repo.Update(ctx, tnt, &spec); err != nil {
		return errors.Wrapf(err, "while updating %s Specification with id %q", spec.ObjectType, spec.ID)
	}

	return nil
}

// UpdateSpecOnlyGlobal takes care of simply updating a single spec entity in db without looking and executing corresponding fetch requests that may be related to it
func (s *service) UpdateSpecOnlyGlobal(ctx context.Context, spec model.Spec) error {
	if err := s.repo.UpdateGlobal(ctx, &spec); err != nil {
		return errors.Wrapf(err, "while updating %s Specification with id %q", spec.ObjectType, spec.ID)
	}

	return nil
}

// DeleteByReferenceObjectID missing godoc
func (s *service) DeleteByReferenceObjectID(ctx context.Context, resourceType resource.Type, objectType model.SpecReferenceObjectType, objectID string) error {
	if err := s.deleteSpecByReferenceObjectID(ctx, objectType, objectID, resourceType); err != nil {
		return errors.Wrapf(err, "while deleting reference object type %s with id %s", objectType, objectID)
	}

	return nil
}

// Delete missing godoc
func (s *service) Delete(ctx context.Context, id string, objectType model.SpecReferenceObjectType) error {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	err = s.repo.Delete(ctx, tnt, id, objectType)
	if err != nil {
		return errors.Wrapf(err, "while deleting Specification with id %q", id)
	}

	return nil
}

// RefetchSpec missing godoc
func (s *service) RefetchSpec(ctx context.Context, id string, objectType model.SpecReferenceObjectType) (*model.Spec, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	spec, err := s.repo.GetByID(ctx, tnt, id, objectType)
	if err != nil {
		return nil, err
	}

	fetchRequest, err := s.fetchRequestRepo.GetByReferenceObjectID(ctx, tnt, getFetchRequestObjectTypeBySpecObjectType(objectType), id)
	if err != nil && !apperrors.IsNotFoundError(err) {
		return nil, errors.Wrapf(err, "while getting FetchRequest for Specification with id %q", id)
	}

	if fetchRequest != nil {
		spec.Data = s.fetchRequestService.HandleSpec(ctx, fetchRequest)
	}

	if err = s.repo.Update(ctx, tnt, spec); err != nil {
		return nil, errors.Wrapf(err, "while updating Specification with id %q", id)
	}

	return spec, nil
}

// GetFetchRequest missing godoc
func (s *service) GetFetchRequest(ctx context.Context, specID string, objectType model.SpecReferenceObjectType) (*model.FetchRequest, error) {
	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.Exists(ctx, tnt, specID, objectType)
	if err != nil {
		return nil, errors.Wrapf(err, "while checking if Specification with id %q exists", specID)
	}
	if !exists {
		return nil, fmt.Errorf("specification with id %q doesn't exist", specID)
	}

	fetchRequest, err := s.fetchRequestRepo.GetByReferenceObjectID(ctx, tnt, getFetchRequestObjectTypeBySpecObjectType(objectType), specID)
	if err != nil {
		return nil, errors.Wrapf(err, "while getting FetchRequest by Specification with id %q", specID)
	}

	return fetchRequest, nil
}

// ListFetchRequestsByReferenceObjectIDs lists specs by reference object IDs
func (s *service) ListFetchRequestsByReferenceObjectIDs(ctx context.Context, tenant string, objectIDs []string, objectType model.SpecReferenceObjectType) ([]*model.FetchRequest, error) {
	return s.fetchRequestRepo.ListByReferenceObjectIDs(ctx, tenant, getFetchRequestObjectTypeBySpecObjectType(objectType), objectIDs)
}

// ListFetchRequestsByReferenceObjectIDsGlobal lists specs by reference object IDs without tenant isolation
func (s *service) ListFetchRequestsByReferenceObjectIDsGlobal(ctx context.Context, objectIDs []string, objectType model.SpecReferenceObjectType) ([]*model.FetchRequest, error) {
	return s.fetchRequestRepo.ListByReferenceObjectIDsGlobal(ctx, getFetchRequestObjectTypeBySpecObjectType(objectType), objectIDs)
}

func (s *service) createFetchRequest(ctx context.Context, in model.FetchRequestInput, parentObjectID string, objectType model.SpecReferenceObjectType, resourceType resource.Type) (*model.FetchRequest, error) {
	id := s.uidService.Generate()
	fr := in.ToFetchRequest(s.timestampGen(), id, getFetchRequestObjectTypeBySpecObjectType(objectType), parentObjectID)

	if resourceType.IsTenantIgnorable() {
		if err := s.fetchRequestRepo.CreateGlobal(ctx, fr); err != nil {
			return nil, err
		}

		return fr, nil
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err = s.fetchRequestRepo.Create(ctx, tnt, fr); err != nil {
		return nil, err
	}

	return fr, nil
}

func (s *service) listSpecsByReferenceObjectID(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string, resourceType resource.Type) ([]*model.Spec, error) {
	if resourceType.IsTenantIgnorable() {
		return s.repo.ListByReferenceObjectIDGlobal(ctx, objectType, objectID)
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.repo.ListByReferenceObjectID(ctx, tnt, objectType, objectID)
}

func (s *service) createSpec(ctx context.Context, spec *model.Spec, resourceType resource.Type) error {
	if resourceType.IsTenantIgnorable() {
		return s.repo.CreateGlobal(ctx, spec)
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repo.Create(ctx, tnt, spec)
}

func (s *service) updateSpec(ctx context.Context, spec *model.Spec, resourceType resource.Type) error {
	if resourceType.IsTenantIgnorable() {
		return s.repo.UpdateGlobal(ctx, spec)
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, tnt, spec)
}

func (s *service) deleteSpecByReferenceObjectID(ctx context.Context, objectType model.SpecReferenceObjectType, objectID string, resourceType resource.Type) error {
	if resourceType.IsTenantIgnorable() {
		return s.repo.DeleteByReferenceObjectIDGlobal(ctx, objectType, objectID)
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	return s.repo.DeleteByReferenceObjectID(ctx, tnt, objectType, objectID)
}

func (s *service) deleteFetchRequestByReferenceObjectID(ctx context.Context, id string, objectType model.SpecReferenceObjectType, resourceType resource.Type) error {
	if resourceType.IsTenantIgnorable() {
		if _, err := s.repo.GetByIDGlobal(ctx, id); err != nil {
			return err
		}

		return s.fetchRequestRepo.DeleteByReferenceObjectIDGlobal(ctx, getFetchRequestObjectTypeBySpecObjectType(objectType), id)
	}

	tnt, err := tenant.LoadFromContext(ctx)
	if err != nil {
		return err
	}

	if _, err = s.repo.GetByID(ctx, tnt, id, objectType); err != nil {
		return err
	}

	return s.fetchRequestRepo.DeleteByReferenceObjectID(ctx, tnt, getFetchRequestObjectTypeBySpecObjectType(objectType), id)
}

func getFetchRequestObjectTypeBySpecObjectType(specObjectType model.SpecReferenceObjectType) model.FetchRequestReferenceObjectType {
	switch specObjectType {
	case model.APISpecReference:
		return model.APISpecFetchRequestReference
	case model.EventSpecReference:
		return model.EventSpecFetchRequestReference
	}
	return ""
}
