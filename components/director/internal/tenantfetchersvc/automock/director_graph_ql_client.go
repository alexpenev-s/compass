// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"
)

// DirectorGraphQLClient is an autogenerated mock type for the DirectorGraphQLClient type
type DirectorGraphQLClient struct {
	mock.Mock
}

// DeleteTenants provides a mock function with given fields: ctx, tenants
func (_m *DirectorGraphQLClient) DeleteTenants(ctx context.Context, tenants []graphql.BusinessTenantMappingInput) error {
	ret := _m.Called(ctx, tenants)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []graphql.BusinessTenantMappingInput) error); ok {
		r0 = rf(ctx, tenants)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubscribeTenant provides a mock function with given fields: ctx, providerID, subaccountID, providerSubaccountID, consumerTenantID, region, subscriptionProviderAppName, subscriptionPayload
func (_m *DirectorGraphQLClient) SubscribeTenant(ctx context.Context, providerID string, subaccountID string, providerSubaccountID string, consumerTenantID string, region string, subscriptionProviderAppName string, subscriptionPayload string) error {
	ret := _m.Called(ctx, providerID, subaccountID, providerSubaccountID, consumerTenantID, region, subscriptionProviderAppName, subscriptionPayload)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, string, string, string) error); ok {
		r0 = rf(ctx, providerID, subaccountID, providerSubaccountID, consumerTenantID, region, subscriptionProviderAppName, subscriptionPayload)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnsubscribeTenant provides a mock function with given fields: ctx, providerID, subaccountID, providerSubaccountID, consumerTenantID, region
func (_m *DirectorGraphQLClient) UnsubscribeTenant(ctx context.Context, providerID string, subaccountID string, providerSubaccountID string, consumerTenantID string, region string) error {
	ret := _m.Called(ctx, providerID, subaccountID, providerSubaccountID, consumerTenantID, region)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, string) error); ok {
		r0 = rf(ctx, providerID, subaccountID, providerSubaccountID, consumerTenantID, region)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTenant provides a mock function with given fields: ctx, id, tenant
func (_m *DirectorGraphQLClient) UpdateTenant(ctx context.Context, id string, tenant graphql.BusinessTenantMappingInput) error {
	ret := _m.Called(ctx, id, tenant)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.BusinessTenantMappingInput) error); ok {
		r0 = rf(ctx, id, tenant)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteTenants provides a mock function with given fields: _a0, _a1
func (_m *DirectorGraphQLClient) WriteTenants(_a0 context.Context, _a1 []graphql.BusinessTenantMappingInput) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []graphql.BusinessTenantMappingInput) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewDirectorGraphQLClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewDirectorGraphQLClient creates a new instance of DirectorGraphQLClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDirectorGraphQLClient(t mockConstructorTestingTNewDirectorGraphQLClient) *DirectorGraphQLClient {
	mock := &DirectorGraphQLClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
