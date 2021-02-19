// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	apptemplate "github.com/kyma-incubator/compass/components/director/internal2/domain/apptemplate"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal2/model"
)

// EntityConverter is an autogenerated mock type for the EntityConverter type
type EntityConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: entity
func (_m *EntityConverter) FromEntity(entity *apptemplate.Entity) (*model.ApplicationTemplate, error) {
	ret := _m.Called(entity)

	var r0 *model.ApplicationTemplate
	if rf, ok := ret.Get(0).(func(*apptemplate.Entity) *model.ApplicationTemplate); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*apptemplate.Entity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *EntityConverter) ToEntity(in *model.ApplicationTemplate) (*apptemplate.Entity, error) {
	ret := _m.Called(in)

	var r0 *apptemplate.Entity
	if rf, ok := ret.Get(0).(func(*model.ApplicationTemplate) *apptemplate.Entity); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*apptemplate.Entity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.ApplicationTemplate) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
