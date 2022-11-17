// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	tenantfetchersvc "github.com/kyma-incubator/compass/components/director/internal/tenantfetchersvc"
	mock "github.com/stretchr/testify/mock"
)

// TenantProvisioner is an autogenerated mock type for the TenantProvisioner type
type TenantProvisioner struct {
	mock.Mock
}

// ProvisionTenants provides a mock function with given fields: _a0, _a1
func (_m *TenantProvisioner) ProvisionTenants(_a0 context.Context, _a1 *tenantfetchersvc.TenantSubscriptionRequest) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *tenantfetchersvc.TenantSubscriptionRequest) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTenantProvisioner interface {
	mock.TestingT
	Cleanup(func())
}

// NewTenantProvisioner creates a new instance of TenantProvisioner. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTenantProvisioner(t mockConstructorTestingTNewTenantProvisioner) *TenantProvisioner {
	mock := &TenantProvisioner{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
