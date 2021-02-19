// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal2/model"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, ld
func (_m *Service) Create(ctx context.Context, ld model.LabelDefinition) (model.LabelDefinition, error) {
	ret := _m.Called(ctx, ld)

	var r0 model.LabelDefinition
	if rf, ok := ret.Get(0).(func(context.Context, model.LabelDefinition) model.LabelDefinition); ok {
		r0 = rf(ctx, ld)
	} else {
		r0 = ret.Get(0).(model.LabelDefinition)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.LabelDefinition) error); ok {
		r1 = rf(ctx, ld)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, tenant, key, deleteRelatedLabels
func (_m *Service) Delete(ctx context.Context, tenant string, key string, deleteRelatedLabels bool) error {
	ret := _m.Called(ctx, tenant, key, deleteRelatedLabels)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, bool) error); ok {
		r0 = rf(ctx, tenant, key, deleteRelatedLabels)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, tenant, key
func (_m *Service) Get(ctx context.Context, tenant string, key string) (*model.LabelDefinition, error) {
	ret := _m.Called(ctx, tenant, key)

	var r0 *model.LabelDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.LabelDefinition); ok {
		r0 = rf(ctx, tenant, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.LabelDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, tenant
func (_m *Service) List(ctx context.Context, tenant string) ([]model.LabelDefinition, error) {
	ret := _m.Called(ctx, tenant)

	var r0 []model.LabelDefinition
	if rf, ok := ret.Get(0).(func(context.Context, string) []model.LabelDefinition); ok {
		r0 = rf(ctx, tenant)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.LabelDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tenant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, ld
func (_m *Service) Update(ctx context.Context, ld model.LabelDefinition) error {
	ret := _m.Called(ctx, ld)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.LabelDefinition) error); ok {
		r0 = rf(ctx, ld)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
