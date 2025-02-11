// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// AggregationFailurePusher is an autogenerated mock type for the AggregationFailurePusher type
type AggregationFailurePusher struct {
	mock.Mock
}

// ReportAggregationFailure provides a mock function with given fields: ctx, err
func (_m *AggregationFailurePusher) ReportAggregationFailure(ctx context.Context, err error) {
	_m.Called(ctx, err)
}

type mockConstructorTestingTNewAggregationFailurePusher interface {
	mock.TestingT
	Cleanup(func())
}

// NewAggregationFailurePusher creates a new instance of AggregationFailurePusher. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewAggregationFailurePusher(t mockConstructorTestingTNewAggregationFailurePusher) *AggregationFailurePusher {
	mock := &AggregationFailurePusher{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
