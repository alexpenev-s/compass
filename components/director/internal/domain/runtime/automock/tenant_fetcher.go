// Code generated by mockery v2.12.2. DO NOT EDIT.

package automock

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// TenantFetcher is an autogenerated mock type for the TenantFetcher type
type TenantFetcher struct {
	mock.Mock
}

// FetchOnDemand provides a mock function with given fields: tenant
func (_m *TenantFetcher) FetchOnDemand(tenant string) error {
	ret := _m.Called(tenant)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(tenant)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTenantFetcher creates a new instance of TenantFetcher. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTenantFetcher(t testing.TB) *TenantFetcher {
	mock := &TenantFetcher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
