// Code generated by mockery. DO NOT EDIT.

package automock

import (
	bundle "github.com/kyma-incubator/compass/components/director/internal/domain/bundle"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	testing "testing"
)

// EntityConverter is an autogenerated mock type for the EntityConverter type
type EntityConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: entity
func (_m *EntityConverter) FromEntity(entity *bundle.Entity) (*model.Bundle, error) {
	ret := _m.Called(entity)

	var r0 *model.Bundle
	if rf, ok := ret.Get(0).(func(*bundle.Entity) *model.Bundle); ok {
		r0 = rf(entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Bundle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*bundle.Entity) error); ok {
		r1 = rf(entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *EntityConverter) ToEntity(in *model.Bundle) (*bundle.Entity, error) {
	ret := _m.Called(in)

	var r0 *bundle.Entity
	if rf, ok := ret.Get(0).(func(*model.Bundle) *bundle.Entity); ok {
		r0 = rf(in)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bundle.Entity)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.Bundle) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewEntityConverter creates a new instance of EntityConverter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewEntityConverter(t testing.TB) *EntityConverter {
	mock := &EntityConverter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
