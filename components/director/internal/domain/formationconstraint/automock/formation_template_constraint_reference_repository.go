// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// FormationTemplateConstraintReferenceRepository is an autogenerated mock type for the formationTemplateConstraintReferenceRepository type
type FormationTemplateConstraintReferenceRepository struct {
	mock.Mock
}

// ListByFormationTemplateID provides a mock function with given fields: ctx, formationTemplateID
func (_m *FormationTemplateConstraintReferenceRepository) ListByFormationTemplateID(ctx context.Context, formationTemplateID string) ([]*model.FormationTemplateConstraintReference, error) {
	ret := _m.Called(ctx, formationTemplateID)

	var r0 []*model.FormationTemplateConstraintReference
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.FormationTemplateConstraintReference); ok {
		r0 = rf(ctx, formationTemplateID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FormationTemplateConstraintReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, formationTemplateID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByFormationTemplateIDs provides a mock function with given fields: ctx, formationTemplateIDs
func (_m *FormationTemplateConstraintReferenceRepository) ListByFormationTemplateIDs(ctx context.Context, formationTemplateIDs []string) ([]*model.FormationTemplateConstraintReference, error) {
	ret := _m.Called(ctx, formationTemplateIDs)

	var r0 []*model.FormationTemplateConstraintReference
	if rf, ok := ret.Get(0).(func(context.Context, []string) []*model.FormationTemplateConstraintReference); ok {
		r0 = rf(ctx, formationTemplateIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.FormationTemplateConstraintReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, formationTemplateIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFormationTemplateConstraintReferenceRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewFormationTemplateConstraintReferenceRepository creates a new instance of FormationTemplateConstraintReferenceRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFormationTemplateConstraintReferenceRepository(t mockConstructorTestingTNewFormationTemplateConstraintReferenceRepository) *FormationTemplateConstraintReferenceRepository {
	mock := &FormationTemplateConstraintReferenceRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
