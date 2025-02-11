// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// SubscriptionService is an autogenerated mock type for the SubscriptionService type
type SubscriptionService struct {
	mock.Mock
}

// SubscribeTenantToRuntime provides a mock function with given fields: ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionAppName, subscriptionID
func (_m *SubscriptionService) SubscribeTenantToRuntime(ctx context.Context, providerID string, subaccountTenantID string, providerSubaccountID string, consumerTenantID string, region string, subscriptionAppName string, subscriptionID string) (bool, error) {
	ret := _m.Called(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionAppName, subscriptionID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, string, string, string) bool); ok {
		r0 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionAppName, subscriptionID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, string, string, string) error); ok {
		r1 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionAppName, subscriptionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UnsubscribeTenantFromRuntime provides a mock function with given fields: ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionID
func (_m *SubscriptionService) UnsubscribeTenantFromRuntime(ctx context.Context, providerID string, subaccountTenantID string, providerSubaccountID string, consumerTenantID string, region string, subscriptionID string) (bool, error) {
	ret := _m.Called(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionID)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string, string, string) bool); ok {
		r0 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string, string, string) error); ok {
		r1 = rf(ctx, providerID, subaccountTenantID, providerSubaccountID, consumerTenantID, region, subscriptionID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSubscriptionService interface {
	mock.TestingT
	Cleanup(func())
}

// NewSubscriptionService creates a new instance of SubscriptionService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSubscriptionService(t mockConstructorTestingTNewSubscriptionService) *SubscriptionService {
	mock := &SubscriptionService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
