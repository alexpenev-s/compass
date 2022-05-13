// Code generated by mockery v2.12.2. DO NOT EDIT.

package automock

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// ScopesGetter is an autogenerated mock type for the ScopesGetter type
type ScopesGetter struct {
	mock.Mock
}

// GetRequiredScopes provides a mock function with given fields: scopesDefinition
func (_m *ScopesGetter) GetRequiredScopes(scopesDefinition string) ([]string, error) {
	ret := _m.Called(scopesDefinition)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(scopesDefinition)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(scopesDefinition)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewScopesGetter creates a new instance of ScopesGetter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewScopesGetter(t testing.TB) *ScopesGetter {
	mock := &ScopesGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
