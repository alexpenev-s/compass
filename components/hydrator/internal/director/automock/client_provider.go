// Code generated by mockery. DO NOT EDIT.

package automock

import (
	director "github.com/kyma-incubator/compass/components/hydrator/internal/director"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// ClientProvider is an autogenerated mock type for the ClientProvider type
type ClientProvider struct {
	mock.Mock
}

// Client provides a mock function with given fields:
func (_m *ClientProvider) Client() director.Client {
	ret := _m.Called()

	var r0 director.Client
	if rf, ok := ret.Get(0).(func() director.Client); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(director.Client)
		}
	}

	return r0
}

// NewClientProvider creates a new instance of ClientProvider. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewClientProvider(t testing.TB) *ClientProvider {
	mock := &ClientProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
