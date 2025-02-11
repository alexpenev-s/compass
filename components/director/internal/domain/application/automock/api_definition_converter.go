// Code generated by mockery. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// APIDefinitionConverter is an autogenerated mock type for the APIDefinitionConverter type
type APIDefinitionConverter struct {
	mock.Mock
}

// ToGraphQL provides a mock function with given fields: in, spec, bundleRef
func (_m *APIDefinitionConverter) ToGraphQL(in *model.APIDefinition, spec *model.Spec, bundleRef *model.BundleReference) (*graphql.APIDefinition, error) {
	ret := _m.Called(in, spec, bundleRef)

	var r0 *graphql.APIDefinition
	if rf, ok := ret.Get(0).(func(*model.APIDefinition, *model.Spec, *model.BundleReference) *graphql.APIDefinition); ok {
		r0 = rf(in, spec, bundleRef)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.APIDefinition)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.APIDefinition, *model.Spec, *model.BundleReference) error); ok {
		r1 = rf(in, spec, bundleRef)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
