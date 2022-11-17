// Code generated by mockery. DO NOT EDIT.

package automock

import (
	api "github.com/kyma-incubator/compass/components/director/internal/domain/api"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// APIDefinitionConverter is an autogenerated mock type for the APIDefinitionConverter type
type APIDefinitionConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: entity
func (_m *APIDefinitionConverter) FromEntity(entity *api.Entity) *model.APIDefinition {
	ret := _m.Called(entity)

	var r0 *model.APIDefinition
	if rf, ok := ret.Get(0).(func(*api.Entity) *model.APIDefinition); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.APIDefinition)
		}
	}

	return r0
}

// ToEntity provides a mock function with given fields: apiModel
func (_m *APIDefinitionConverter) ToEntity(apiModel *model.APIDefinition) *api.Entity {
	ret := _m.Called(apiModel)

	var r0 *api.Entity
	if rf, ok := ret.Get(0).(func(*model.APIDefinition) *api.Entity); ok {
		r0 = rf(apiModel)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*api.Entity)
		}
	}

	return r0
}

type mockConstructorTestingTNewAPIDefinitionConverter interface {
	mock.TestingT
	Cleanup(func())
}

// NewAPIDefinitionConverter creates a new instance of APIDefinitionConverter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAPIDefinitionConverter(t mockConstructorTestingTNewAPIDefinitionConverter) *APIDefinitionConverter {
	mock := &APIDefinitionConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
