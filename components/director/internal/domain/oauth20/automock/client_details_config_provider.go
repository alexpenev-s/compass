// Code generated by mockery. DO NOT EDIT.

package automock

import mock "github.com/stretchr/testify/mock"

// ClientDetailsConfigProvider is an autogenerated mock type for the ClientDetailsConfigProvider type
type ClientDetailsConfigProvider struct {
	mock.Mock
}

// GetRequiredGrantTypes provides a mock function with given fields: path
func (_m *ClientDetailsConfigProvider) GetRequiredGrantTypes(path string) ([]string, error) {
	ret := _m.Called(path)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRequiredScopes provides a mock function with given fields: path
func (_m *ClientDetailsConfigProvider) GetRequiredScopes(path string) ([]string, error) {
	ret := _m.Called(path)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewClientDetailsConfigProvider interface {
	mock.TestingT
	Cleanup(func())
}

// NewClientDetailsConfigProvider creates a new instance of ClientDetailsConfigProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClientDetailsConfigProvider(t mockConstructorTestingTNewClientDetailsConfigProvider) *ClientDetailsConfigProvider {
	mock := &ClientDetailsConfigProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
