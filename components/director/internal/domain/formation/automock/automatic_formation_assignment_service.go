// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// AutomaticFormationAssignmentService is an autogenerated mock type for the automaticFormationAssignmentService type
type AutomaticFormationAssignmentService struct {
	mock.Mock
}

// GetForScenarioName provides a mock function with given fields: ctx, scenarioName
func (_m *AutomaticFormationAssignmentService) GetForScenarioName(ctx context.Context, scenarioName string) (model.AutomaticScenarioAssignment, error) {
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

type mockConstructorTestingTNewAutomaticFormationAssignmentService interface {
	mock.TestingT
	Cleanup(func())
}

// NewAutomaticFormationAssignmentService creates a new instance of AutomaticFormationAssignmentService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAutomaticFormationAssignmentService(t mockConstructorTestingTNewAutomaticFormationAssignmentService) *AutomaticFormationAssignmentService {
	mock := &AutomaticFormationAssignmentService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
