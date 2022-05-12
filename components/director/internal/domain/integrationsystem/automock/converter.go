// Code generated by mockery v2.12.1. DO NOT EDIT.

package automock

import (
	integrationsystem "github.com/kyma-incubator/compass/components/director/internal/domain/integrationsystem"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// Converter is an autogenerated mock type for the Converter type
type Converter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: in
func (_m *Converter) FromEntity(in *integrationsystem.Entity) *model.IntegrationSystem {
	ret := _m.Called(in)

	var r0 *model.IntegrationSystem
	if rf, ok := ret.Get(0).(func(*integrationsystem.Entity) *model.IntegrationSystem); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.IntegrationSystem)
		}
	}

	return r0
}

// ToEntity provides a mock function with given fields: in
func (_m *Converter) ToEntity(in *model.IntegrationSystem) *integrationsystem.Entity {
	ret := _m.Called(in)

	var r0 *integrationsystem.Entity
	if rf, ok := ret.Get(0).(func(*model.IntegrationSystem) *integrationsystem.Entity); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*integrationsystem.Entity)
		}
	}

	return r0
}

// NewConverter creates a new instance of Converter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewConverter(t testing.TB) *Converter {
	mock := &Converter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
