// Code generated by mockery. DO NOT EDIT.

package automock

import mock "github.com/stretchr/testify/mock"

// Status is an autogenerated mock type for the Status type
type Status struct {
	mock.Mock
}

// Details provides a mock function with given fields:
func (_m *Status) Details() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Error provides a mock function with given fields:
func (_m *Status) Error() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewStatus interface {
	mock.TestingT
	Cleanup(func())
}

// NewStatus creates a new instance of Status. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewStatus(t mockConstructorTestingTNewStatus) *Status {
	mock := &Status{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
