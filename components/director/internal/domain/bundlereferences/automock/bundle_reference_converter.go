// Code generated by mockery v2.2.1. DO NOT EDIT.

package automock

import (
	bundlereferences "github.com/kyma-incubator/compass/components/director/internal/domain/bundlereferences"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// BundleReferenceConverter is an autogenerated mock type for the BundleReferenceConverter type
type BundleReferenceConverter struct {
	mock.Mock
}

// FromEntity provides a mock function with given fields: in
func (_m *BundleReferenceConverter) FromEntity(in bundlereferences.Entity) (model.BundleReference, error) {
	ret := _m.Called(in)

	var r0 model.BundleReference
	if rf, ok := ret.Get(0).(func(bundlereferences.Entity) model.BundleReference); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(model.BundleReference)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bundlereferences.Entity) error); ok {
		r1 = rf(in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ToEntity provides a mock function with given fields: in
func (_m *BundleReferenceConverter) ToEntity(in model.BundleReference) bundlereferences.Entity {
	ret := _m.Called(in)

	var r0 bundlereferences.Entity
	if rf, ok := ret.Get(0).(func(model.BundleReference) bundlereferences.Entity); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(bundlereferences.Entity)
	}

	return r0
}
