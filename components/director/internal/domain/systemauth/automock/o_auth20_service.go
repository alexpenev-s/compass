// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// OAuth20Service is an autogenerated mock type for the OAuth20Service type
type OAuth20Service struct {
	mock.Mock
}

// DeleteClientCredentials provides a mock function with given fields: ctx, clientID
func (_m *OAuth20Service) DeleteClientCredentials(ctx context.Context, clientID string) error {
	ret := _m.Called(ctx, clientID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, clientID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewOAuth20Service interface {
	mock.TestingT
	Cleanup(func())
}

// NewOAuth20Service creates a new instance of OAuth20Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewOAuth20Service(t mockConstructorTestingTNewOAuth20Service) *OAuth20Service {
	mock := &OAuth20Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
