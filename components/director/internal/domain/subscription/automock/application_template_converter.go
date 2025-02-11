// Code generated by mockery. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// ApplicationTemplateConverter is an autogenerated mock type for the ApplicationTemplateConverter type
type ApplicationTemplateConverter struct {
	mock.Mock
}

// ApplicationFromTemplateInputFromGraphQL provides a mock function with given fields: appTemplate, in
func (_m *ApplicationTemplateConverter) ApplicationFromTemplateInputFromGraphQL(appTemplate *model.ApplicationTemplate, in graphql.ApplicationFromTemplateInput) (model.ApplicationFromTemplateInput, error) {
	ret := _m.Called(appTemplate, in)

	var r0 model.ApplicationFromTemplateInput
	if rf, ok := ret.Get(0).(func(*model.ApplicationTemplate, graphql.ApplicationFromTemplateInput) model.ApplicationFromTemplateInput); ok {
		r0 = rf(appTemplate, in)
	} else {
		r0 = ret.Get(0).(model.ApplicationFromTemplateInput)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.ApplicationTemplate, graphql.ApplicationFromTemplateInput) error); ok {
		r1 = rf(appTemplate, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewApplicationTemplateConverter interface {
	mock.TestingT
	Cleanup(func())
}

// NewApplicationTemplateConverter creates a new instance of ApplicationTemplateConverter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewApplicationTemplateConverter(t mockConstructorTestingTNewApplicationTemplateConverter) *ApplicationTemplateConverter {
	mock := &ApplicationTemplateConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
