// Code generated by mockery. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// Converter is an autogenerated mock type for the Converter type
type Converter struct {
	mock.Mock
}

// FromGraphQL provides a mock function with given fields: i
func (_m *Converter) FromGraphQL(i graphql.FormationInput) model.Formation {
	ret := _m.Called(i)

	var r0 model.Formation
	if rf, ok := ret.Get(0).(func(graphql.FormationInput) model.Formation); ok {
		r0 = rf(i)
	} else {
		r0 = ret.Get(0).(model.Formation)
	}

	return r0
}

// MultipleToGraphQL provides a mock function with given fields: in
func (_m *Converter) MultipleToGraphQL(in []*model.Formation) ([]*graphql.Formation, error) {
	ret := _m.Called(in)

	var r0 []*graphql.Formation
	if rf, ok := ret.Get(0).(func([]*model.Formation) []*graphql.Formation); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*graphql.Formation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]*model.Formation) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToGraphQL provides a mock function with given fields: i
func (_m *Converter) ToGraphQL(i *model.Formation) (*graphql.Formation, error) {
	ret := _m.Called(i)

	var r0 *graphql.Formation
	if rf, ok := ret.Get(0).(func(*model.Formation) *graphql.Formation); ok {
		r0 = rf(i)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*graphql.Formation)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Formation) error); ok {
		r1 = rf(i)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewConverter interface {
	mock.TestingT
	Cleanup(func())
}

// NewConverter creates a new instance of Converter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConverter(t mockConstructorTestingTNewConverter) *Converter {
	mock := &Converter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
