// Code generated by mockery v2.12.2. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// ScenarioAssignmentEngine is an autogenerated mock type for the scenarioAssignmentEngine type
type ScenarioAssignmentEngine struct {
	mock.Mock
}

// MergeScenariosFromInputLabelsAndAssignments provides a mock function with given fields: ctx, inputLabels, runtimeID
func (_m *ScenarioAssignmentEngine) MergeScenariosFromInputLabelsAndAssignments(ctx context.Context, inputLabels map[string]interface{}, runtimeID string) ([]interface{}, error) {
	ret := _m.Called(ctx, inputLabels, runtimeID)

	var r0 []interface{}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]interface{}, string) []interface{}); ok {
		r0 = rf(ctx, inputLabels, runtimeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, map[string]interface{}, string) error); ok {
		r1 = rf(ctx, inputLabels, runtimeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewScenarioAssignmentEngine creates a new instance of ScenarioAssignmentEngine. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewScenarioAssignmentEngine(t testing.TB) *ScenarioAssignmentEngine {
	mock := &ScenarioAssignmentEngine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
