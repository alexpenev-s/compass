// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// TenantService is an autogenerated mock type for the tenantService type
type TenantService struct {
	mock.Mock
}

// ListsByExternalIDs provides a mock function with given fields: ctx, ids
func (_m *TenantService) ListsByExternalIDs(ctx context.Context, ids []string) ([]*model.BusinessTenantMapping, error) {
	ret := _m.Called(ctx, ids)

	var r0 []*model.BusinessTenantMapping
	if rf, ok := ret.Get(0).(func(context.Context, []string) []*model.BusinessTenantMapping); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.BusinessTenantMapping)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, ids)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTenantService interface {
	mock.TestingT
	Cleanup(func())
}

// NewTenantService creates a new instance of TenantService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTenantService(t mockConstructorTestingTNewTenantService) *TenantService {
	mock := &TenantService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
