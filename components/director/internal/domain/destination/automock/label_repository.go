// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// LabelRepository is an autogenerated mock type for the labelRepository type
type LabelRepository struct {
	mock.Mock
}

// GetByKey provides a mock function with given fields: ctx, tenant, objectType, objectID, key
func (_m *LabelRepository) GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, key string) (*model.Label, error) {
	ret := _m.Called(ctx, tenant, objectType, objectID, key)

	var r0 *model.Label
	if rf, ok := ret.Get(0).(func(context.Context, string, model.LabelableObject, string, string) *model.Label); ok {
		r0 = rf(ctx, tenant, objectType, objectID, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.LabelableObject, string, string) error); ok {
		r1 = rf(ctx, tenant, objectType, objectID, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListForGlobalObject provides a mock function with given fields: ctx, objectType, objectID
func (_m *LabelRepository) ListForGlobalObject(ctx context.Context, objectType model.LabelableObject, objectID string) (map[string]*model.Label, error) {
	ret := _m.Called(ctx, objectType, objectID)

	var r0 map[string]*model.Label
	if rf, ok := ret.Get(0).(func(context.Context, model.LabelableObject, string) map[string]*model.Label); ok {
		r0 = rf(ctx, objectType, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.LabelableObject, string) error); ok {
		r1 = rf(ctx, objectType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListForObject provides a mock function with given fields: ctx, tenant, objectType, objectID
func (_m *LabelRepository) ListForObject(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string) (map[string]*model.Label, error) {
	ret := _m.Called(ctx, tenant, objectType, objectID)

	var r0 map[string]*model.Label
	if rf, ok := ret.Get(0).(func(context.Context, string, model.LabelableObject, string) map[string]*model.Label); ok {
		r0 = rf(ctx, tenant, objectType, objectID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.LabelableObject, string) error); ok {
		r1 = rf(ctx, tenant, objectType, objectID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewLabelRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewLabelRepository creates a new instance of LabelRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLabelRepository(t mockConstructorTestingTNewLabelRepository) *LabelRepository {
	mock := &LabelRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
