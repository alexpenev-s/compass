// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// TenantRepo is an autogenerated mock type for the TenantRepo type
type TenantRepo struct {
	mock.Mock
}

// ListBySubscribedRuntimes provides a mock function with given fields: ctx
func (_m *TenantRepo) ListBySubscribedRuntimes(ctx context.Context) ([]*model.BusinessTenantMapping, error) {
	ret := _m.Called(ctx)

	var r0 []*model.BusinessTenantMapping
	if rf, ok := ret.Get(0).(func(context.Context) []*model.BusinessTenantMapping); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.BusinessTenantMapping)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTenantRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewTenantRepo creates a new instance of TenantRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTenantRepo(t mockConstructorTestingTNewTenantRepo) *TenantRepo {
	mock := &TenantRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
