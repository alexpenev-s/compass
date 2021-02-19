// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal2/model"
	mock "github.com/stretchr/testify/mock"
)

// AssignmentEngine is an autogenerated mock type for the AssignmentEngine type
type AssignmentEngine struct {
	mock.Mock
}

// EnsureScenarioAssigned provides a mock function with given fields: ctx, in
func (_m *AssignmentEngine) EnsureScenarioAssigned(ctx context.Context, in model.AutomaticScenarioAssignment) error {
	ret := _m.Called(ctx, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.AutomaticScenarioAssignment) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveAssignedScenario provides a mock function with given fields: ctx, in
func (_m *AssignmentEngine) RemoveAssignedScenario(ctx context.Context, in model.AutomaticScenarioAssignment) error {
	ret := _m.Called(ctx, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.AutomaticScenarioAssignment) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveAssignedScenarios provides a mock function with given fields: ctx, in
func (_m *AssignmentEngine) RemoveAssignedScenarios(ctx context.Context, in []*model.AutomaticScenarioAssignment) error {
	ret := _m.Called(ctx, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*model.AutomaticScenarioAssignment) error); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
