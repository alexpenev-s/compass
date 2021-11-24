// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	graphql "github.com/kyma-incubator/compass/components/director/pkg/graphql"
	mock "github.com/stretchr/testify/mock"
)

// DirectorGraphQLClient is an autogenerated mock type for the DirectorGraphQLClient type
type DirectorGraphQLClient struct {
	mock.Mock
}

// CreateLabelDefinition provides a mock function with given fields: _a0, _a1, _a2
func (_m *DirectorGraphQLClient) CreateLabelDefinition(_a0 context.Context, _a1 graphql.LabelDefinitionInput, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, graphql.LabelDefinitionInput, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTenants provides a mock function with given fields: _a0, _a1
func (_m *DirectorGraphQLClient) DeleteTenants(_a0 context.Context, _a1 []graphql.BusinessTenantMappingInput) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []graphql.BusinessTenantMappingInput) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetRuntimeTenant provides a mock function with given fields: ctx, runtimeID, tenantID, tenantHeader
func (_m *DirectorGraphQLClient) SetRuntimeTenant(ctx context.Context, runtimeID string, tenantID string, tenantHeader string) error {
	ret := _m.Called(ctx, runtimeID, tenantID, tenantHeader)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) error); ok {
		r0 = rf(ctx, runtimeID, tenantID, tenantHeader)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateLabelDefinition provides a mock function with given fields: _a0, _a1, _a2
func (_m *DirectorGraphQLClient) UpdateLabelDefinition(_a0 context.Context, _a1 graphql.LabelDefinitionInput, _a2 string) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, graphql.LabelDefinitionInput, string) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTenant provides a mock function with given fields: _a0, _a1, _a2
func (_m *DirectorGraphQLClient) UpdateTenant(_a0 context.Context, _a1 string, _a2 graphql.BusinessTenantMappingInput) error {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, graphql.BusinessTenantMappingInput) error); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WriteTenants provides a mock function with given fields: _a0, _a1
func (_m *DirectorGraphQLClient) WriteTenants(_a0 context.Context, _a1 []graphql.BusinessTenantMappingInput) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []graphql.BusinessTenantMappingInput) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
