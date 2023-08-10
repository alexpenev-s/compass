// Code generated by mockery. DO NOT EDIT.

package automock

import (
	formationtemplate "github.com/kyma-incubator/compass/components/director/internal/domain/formationtemplate"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// EntityConverter is an autogenerated mock type for the EntityConverter type
type EntityConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: entity
func (_m *EntityConverter) FromEntity(entity *formationtemplate.Entity) (*model.FormationTemplate, error) {
	ret := _m.Called(entity)

	var r0 *model.FormationTemplate
	var r1 error
	if rf, ok := ret.Get(0).(func(*formationtemplate.Entity) (*model.FormationTemplate, error)); ok {
		return rf(entity)
	}
	if rf, ok := ret.Get(0).(func(*formationtemplate.Entity) *model.FormationTemplate); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationTemplate)
		}
	}

	if rf, ok := ret.Get(1).(func(*formationtemplate.Entity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *EntityConverter) ToEntity(in *model.FormationTemplate) (*formationtemplate.Entity, error) {
	ret := _m.Called(in)

	var r0 *formationtemplate.Entity
	var r1 error
	if rf, ok := ret.Get(0).(func(*model.FormationTemplate) (*formationtemplate.Entity, error)); ok {
		return rf(in)
	}
	if rf, ok := ret.Get(0).(func(*model.FormationTemplate) *formationtemplate.Entity); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*formationtemplate.Entity)
		}
	}

	if rf, ok := ret.Get(1).(func(*model.FormationTemplate) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewEntityConverter interface {
	mock.TestingT
	Cleanup(func())
}

// NewEntityConverter creates a new instance of EntityConverter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEntityConverter(t mockConstructorTestingTNewEntityConverter) *EntityConverter {
	mock := &EntityConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
