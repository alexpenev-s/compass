// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// LabelRepository is an autogenerated mock type for the labelRepository type
type LabelRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, tenant, objectType, objectID, key
func (_m *LabelRepository) Delete(ctx context.Context, tenant string, objectType model.LabelableObject, objectID string, key string) error {
	ret := _m.Called(ctx, tenant, objectType, objectID, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.LabelableObject, string, string) error); ok {
		r0 = rf(ctx, tenant, objectType, objectID, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetScenarioLabelsForRuntimes provides a mock function with given fields: ctx, tenantID, runtimesIDs
func (_m *LabelRepository) GetScenarioLabelsForRuntimes(ctx context.Context, tenantID string, runtimesIDs []string) ([]model.Label, error) {
	ret := _m.Called(ctx, tenantID, runtimesIDs)

	var r0 []model.Label
	if rf, ok := ret.Get(0).(func(context.Context, string, []string) []model.Label); ok {
		r0 = rf(ctx, tenantID, runtimesIDs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []string) error); ok {
		r1 = rf(ctx, tenantID, runtimesIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
