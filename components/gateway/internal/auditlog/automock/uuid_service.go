// Code generated by mockery. DO NOT EDIT.

package automock

import mock "github.com/stretchr/testify/mock"

// UUIDService is an autogenerated mock type for the UUIDService type
type UUIDService struct {
	mock.Mock
}

// Generate provides a mock function with given fields:
func (_m *UUIDService) Generate() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
