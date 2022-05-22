// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// ScenarioAssignmentService is an autogenerated mock type for the ScenarioAssignmentService type
type ScenarioAssignmentService struct {
	mock.Mock
}

// Delete provides a mock function with given fields: ctx, in
func (_m *ScenarioAssignmentService) Delete(ctx context.Context, in model.AutomaticScenarioAssignment) error {
	ret := _m.Called(ctx, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.AutomaticScenarioAssignment) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetForScenarioName provides a mock function with given fields: ctx, scenarioName
func (_m *ScenarioAssignmentService) GetForScenarioName(ctx context.Context, scenarioName string) (model.AutomaticScenarioAssignment, error) {
	ret := _m.Called(ctx, scenarioName)

	var r0 model.AutomaticScenarioAssignment
	if rf, ok := ret.Get(0).(func(context.Context, string) model.AutomaticScenarioAssignment); ok {
		r0 = rf(ctx, scenarioName)
	} else {
		r0 = ret.Get(0).(model.AutomaticScenarioAssignment)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, scenarioName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
