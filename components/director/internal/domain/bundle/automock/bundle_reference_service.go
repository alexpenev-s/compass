// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// BundleReferenceService is an autogenerated mock type for the BundleReferenceService type
type BundleReferenceService struct {
	mock.Mock
}

// GetForBundle provides a mock function with given fields: ctx, objectType, objectID, bundleID
func (_m *BundleReferenceService) GetForBundle(ctx context.Context, objectType model.BundleReferenceObjectType, objectID *string, bundleID *string) (*model.BundleReference, error) {
	ret := _m.Called(ctx, objectType, objectID, bundleID)

	var r0 *model.BundleReference
	if rf, ok := ret.Get(0).(func(context.Context, model.BundleReferenceObjectType, *string, *string) *model.BundleReference); ok {
		r0 = rf(ctx, objectType, objectID, bundleID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.BundleReference)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.BundleReferenceObjectType, *string, *string) error); ok {
		r1 = rf(ctx, objectType, objectID, bundleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByBundleIDs provides a mock function with given fields: ctx, objectType, bundleIDs, pageSize, cursor
func (_m *BundleReferenceService) ListByBundleIDs(ctx context.Context, objectType model.BundleReferenceObjectType, bundleIDs []string, pageSize int, cursor string) ([]*model.BundleReference, map[string]int, error) {
	ret := _m.Called(ctx, objectType, bundleIDs, pageSize, cursor)

	var r0 []*model.BundleReference
	if rf, ok := ret.Get(0).(func(context.Context, model.BundleReferenceObjectType, []string, int, string) []*model.BundleReference); ok {
		r0 = rf(ctx, objectType, bundleIDs, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.BundleReference)
		}
	}

	var r1 map[string]int
	if rf, ok := ret.Get(1).(func(context.Context, model.BundleReferenceObjectType, []string, int, string) map[string]int); ok {
		r1 = rf(ctx, objectType, bundleIDs, pageSize, cursor)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(map[string]int)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, model.BundleReferenceObjectType, []string, int, string) error); ok {
		r2 = rf(ctx, objectType, bundleIDs, pageSize, cursor)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewBundleReferenceService interface {
	mock.TestingT
	Cleanup(func())
}

// NewBundleReferenceService creates a new instance of BundleReferenceService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBundleReferenceService(t mockConstructorTestingTNewBundleReferenceService) *BundleReferenceService {
	mock := &BundleReferenceService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
