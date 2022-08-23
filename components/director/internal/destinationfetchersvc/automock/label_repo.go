// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// LabelRepo is an autogenerated mock type for the LabelRepo type
type LabelRepo struct {
	mock.Mock
}

// GetByKey provides a mock function with given fields: ctx, tenant, objectType, objectID, key
func (_m *LabelRepo) GetByKey(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, key string) (*model.Label, error) {
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

// GetSubdomainLabelForSubscribedRuntime provides a mock function with given fields: ctx, tenantID
func (_m *LabelRepo) GetSubdomainLabelForSubscribedRuntime(ctx context.Context, tenantID string) (*model.Label, error) {
	ret := _m.Called(ctx, tenantID)

	var r0 *model.Label
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Label); ok {
		r0 = rf(ctx, tenantID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, tenantID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewLabelRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewLabelRepo creates a new instance of LabelRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLabelRepo(t mockConstructorTestingTNewLabelRepo) *LabelRepo {
	mock := &LabelRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
