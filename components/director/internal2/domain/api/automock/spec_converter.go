// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal2/model"
)

// SpecConverter is an autogenerated mock type for the SpecConverter type
type SpecConverter struct {
	mock.Mock
}

// InputFromGraphQLAPISpec provides a mock function with given fields: in
func (_m *SpecConverter) InputFromGraphQLAPISpec(in *graphql.APISpecInput) (*model.SpecInput, error) {
	ret := _m.Called(in)

	var r0 *model.SpecInput
	if rf, ok := ret.Get(0).(func(*graphql.APISpecInput) *model.SpecInput); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SpecInput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*graphql.APISpecInput) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToGraphQLAPISpec provides a mock function with given fields: in
func (_m *SpecConverter) ToGraphQLAPISpec(in *model.Spec) (*graphql.APISpec, error) {
	ret := _m.Called(in)

	var r0 *graphql.APISpec
	if rf, ok := ret.Get(0).(func(*model.Spec) *graphql.APISpec); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.APISpec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Spec) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
