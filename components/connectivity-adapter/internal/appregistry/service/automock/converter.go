// Code generated by mockery v2.5.1. DO NOT EDIT.

package automock

import (
	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/connectivity-adapter/pkg/model"

	service "github.com/kyma-incubator/compass/components/connectivity-adapter/internal/appregistry/service"
)

// Converter is an autogenerated mock type for the Converter type
type Converter struct {
	mock.Mock
}

// DetailsToGraphQLCreateInput provides a mock function with given fields: deprecated
func (_m *Converter) DetailsToGraphQLCreateInput(deprecated model.ServiceDetails) (graphql.BundleCreateInput, error) {
	ret := _m.Called(deprecated)

	var r0 graphql.BundleCreateInput
	if rf, ok := ret.Get(0).(func(model.ServiceDetails) graphql.BundleCreateInput); ok {
		r0 = rf(deprecated)
	} else {
		r0 = ret.Get(0).(graphql.BundleCreateInput)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.ServiceDetails) error); ok {
		r1 = rf(deprecated)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GraphQLCreateInputToUpdateInput provides a mock function with given fields: in
func (_m *Converter) GraphQLCreateInputToUpdateInput(in graphql.BundleCreateInput) graphql.BundleUpdateInput {
	ret := _m.Called(in)

	var r0 graphql.BundleUpdateInput
	if rf, ok := ret.Get(0).(func(graphql.BundleCreateInput) graphql.BundleUpdateInput); ok {
		r0 = rf(in)
	} else {
		r0 = ret.Get(0).(graphql.BundleUpdateInput)
	}

	return r0
}

// GraphQLToServiceDetails provides a mock function with given fields: converted, legacyServiceReference
func (_m *Converter) GraphQLToServiceDetails(converted graphql.BundleExt, legacyServiceReference service.LegacyServiceReference) (model.ServiceDetails, error) {
	ret := _m.Called(converted, legacyServiceReference)

	var r0 model.ServiceDetails
	if rf, ok := ret.Get(0).(func(graphql.BundleExt, service.LegacyServiceReference) model.ServiceDetails); ok {
		r0 = rf(converted, legacyServiceReference)
	} else {
		r0 = ret.Get(0).(model.ServiceDetails)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(graphql.BundleExt, service.LegacyServiceReference) error); ok {
		r1 = rf(converted, legacyServiceReference)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ServiceDetailsToService provides a mock function with given fields: in, serviceID
func (_m *Converter) ServiceDetailsToService(in model.ServiceDetails, serviceID string) (model.Service, error) {
	ret := _m.Called(in, serviceID)

	var r0 model.Service
	if rf, ok := ret.Get(0).(func(model.ServiceDetails, string) model.Service); ok {
		r0 = rf(in, serviceID)
	} else {
		r0 = ret.Get(0).(model.Service)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.ServiceDetails, string) error); ok {
		r1 = rf(in, serviceID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
