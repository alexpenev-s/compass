// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	tenantmapping "github.com/kyma-incubator/compass/components/director/internal2/tenantmapping"
	mock "github.com/stretchr/testify/mock"
)

// StaticUserRepository is an autogenerated mock type for the StaticUserRepository type
type StaticUserRepository struct {
	mock.Mock
}

// Get provides a mock function with given fields: username
func (_m *StaticUserRepository) Get(username string) (tenantmapping.StaticUser, error) {
	ret := _m.Called(username)

	var r0 tenantmapping.StaticUser
	if rf, ok := ret.Get(0).(func(string) tenantmapping.StaticUser); ok {
		r0 = rf(username)
	} else {
		r0 = ret.Get(0).(tenantmapping.StaticUser)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(username)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
