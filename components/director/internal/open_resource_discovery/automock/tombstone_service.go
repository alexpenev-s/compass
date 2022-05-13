// Code generated by mockery v2.12.2. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// TombstoneService is an autogenerated mock type for the TombstoneService type
type TombstoneService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, applicationID, in
func (_m *TombstoneService) Create(ctx context.Context, applicationID string, in model.TombstoneInput) (string, error) {
	ret := _m.Called(ctx, applicationID, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, model.TombstoneInput) string); ok {
		r0 = rf(ctx, applicationID, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.TombstoneInput) error); ok {
		r1 = rf(ctx, applicationID, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *TombstoneService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ListByApplicationID provides a mock function with given fields: ctx, appID
func (_m *TombstoneService) ListByApplicationID(ctx context.Context, appID string) ([]*model.Tombstone, error) {
	ret := _m.Called(ctx, appID)

	var r0 []*model.Tombstone
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Tombstone); ok {
		r0 = rf(ctx, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Tombstone)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, in
func (_m *TombstoneService) Update(ctx context.Context, id string, in model.TombstoneInput) error {
	ret := _m.Called(ctx, id, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.TombstoneInput) error); ok {
		r0 = rf(ctx, id, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTombstoneService creates a new instance of TombstoneService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTombstoneService(t testing.TB) *TombstoneService {
	mock := &TombstoneService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
