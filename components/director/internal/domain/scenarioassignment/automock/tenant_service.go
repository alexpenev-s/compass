// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TenantService is an autogenerated mock type for the tenantService type
type TenantService struct {
	mock.Mock
}

// GetExternalTenant provides a mock function with given fields: ctx, id
func (_m *TenantService) GetExternalTenant(ctx context.Context, id string) (string, error) {
	ret := _m.Called(ctx, id)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetInternalTenant provides a mock function with given fields: ctx, externalTenant
func (_m *TenantService) GetInternalTenant(ctx context.Context, externalTenant string) (string, error) {
	ret := _m.Called(ctx, externalTenant)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, externalTenant)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, externalTenant)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type NewTenantServiceT interface {
	mock.TestingT
	Cleanup(func())
}

// NewTenantService creates a new instance of TenantService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTenantService(t NewTenantServiceT) *TenantService {
	mock := &TenantService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
