// Code generated by mockery. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// FormationConstraintConverter is an autogenerated mock type for the formationConstraintConverter type
type FormationConstraintConverter struct {
	mock.Mock
}

// FromInputGraphQL provides a mock function with given fields: in
func (_m *FormationConstraintConverter) FromInputGraphQL(in *graphql.FormationConstraintInput) *model.FormationConstraintInput {
	ret := _m.Called(in)

	var r0 *model.FormationConstraintInput
	if rf, ok := ret.Get(0).(func(*graphql.FormationConstraintInput) *model.FormationConstraintInput); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationConstraintInput)
		}
	}

	return r0
}

// FromModelInputToModel provides a mock function with given fields: in, id
func (_m *FormationConstraintConverter) FromModelInputToModel(in *model.FormationConstraintInput, id string) *model.FormationConstraint {
	ret := _m.Called(in, id)

	var r0 *model.FormationConstraint
	if rf, ok := ret.Get(0).(func(*model.FormationConstraintInput, string) *model.FormationConstraint); ok {
		r0 = rf(in, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationConstraint)
		}
	}

	return r0
}

// MultipleToGraphQL provides a mock function with given fields: in
func (_m *FormationConstraintConverter) MultipleToGraphQL(in []*model.FormationConstraint) []*graphql.FormationConstraint {
	ret := _m.Called(in)

	var r0 []*graphql.FormationConstraint
	if rf, ok := ret.Get(0).(func([]*model.FormationConstraint) []*graphql.FormationConstraint); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.FormationConstraint)
		}
	}

	return r0
}

// ToGraphQL provides a mock function with given fields: in
func (_m *FormationConstraintConverter) ToGraphQL(in *model.FormationConstraint) *graphql.FormationConstraint {
	ret := _m.Called(in)

	var r0 *graphql.FormationConstraint
	if rf, ok := ret.Get(0).(func(*model.FormationConstraint) *graphql.FormationConstraint); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.FormationConstraint)
		}
	}

	return r0
}

type mockConstructorTestingTNewFormationConstraintConverter interface {
	mock.TestingT
	Cleanup(func())
}

// NewFormationConstraintConverter creates a new instance of FormationConstraintConverter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFormationConstraintConverter(t mockConstructorTestingTNewFormationConstraintConverter) *FormationConstraintConverter {
	mock := &FormationConstraintConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
