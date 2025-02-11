// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	resource "github.com/kyma-incubator/compass/components/director/pkg/resource"
)

// TombstoneService is an autogenerated mock type for the TombstoneService type
type TombstoneService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, resourceType, resourceID, in
func (_m *TombstoneService) Create(ctx context.Context, resourceType resource.Type, resourceID string, in model.TombstoneInput) (string, error) {
	ret := _m.Called(ctx, resourceType, resourceID, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, resource.Type, string, model.TombstoneInput) string); ok {
		r0 = rf(ctx, resourceType, resourceID, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, resource.Type, string, model.TombstoneInput) error); ok {
		r1 = rf(ctx, resourceType, resourceID, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByApplicationID provides a mock function with given fields: ctx, appID
func (_m *TombstoneService) ListByApplicationID(ctx context.Context, appID string) ([]*model.Tombstone, error) {
	ret := _m.Called(ctx, appID)

	var r0 []*model.Tombstone
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Tombstone); ok {
		r0 = rf(ctx, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Tombstone)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByApplicationTemplateVersionID provides a mock function with given fields: ctx, appID
func (_m *TombstoneService) ListByApplicationTemplateVersionID(ctx context.Context, appID string) ([]*model.Tombstone, error) {
	ret := _m.Called(ctx, appID)

	var r0 []*model.Tombstone
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Tombstone); ok {
		r0 = rf(ctx, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Tombstone)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, resourceType, id, in
func (_m *TombstoneService) Update(ctx context.Context, resourceType resource.Type, id string, in model.TombstoneInput) error {
	ret := _m.Called(ctx, resourceType, id, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, resource.Type, string, model.TombstoneInput) error); ok {
		r0 = rf(ctx, resourceType, id, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTombstoneService interface {
	mock.TestingT
	Cleanup(func())
}

// NewTombstoneService creates a new instance of TombstoneService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTombstoneService(t mockConstructorTestingTNewTombstoneService) *TombstoneService {
	mock := &TombstoneService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
