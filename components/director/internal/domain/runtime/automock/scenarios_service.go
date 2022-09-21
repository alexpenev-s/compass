// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ScenariosService is an autogenerated mock type for the scenariosService type
type ScenariosService struct {
	mock.Mock
}

// EnsureScenariosLabelDefinitionExists provides a mock function with given fields: ctx, tenant
func (_m *ScenariosService) EnsureScenariosLabelDefinitionExists(ctx context.Context, tenant string) error {
	ret := _m.Called(ctx, tenant)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, tenant)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
