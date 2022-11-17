// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// BundleRepository is an autogenerated mock type for the BundleRepository type
type BundleRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, tenant, item
func (_m *BundleRepository) Create(ctx context.Context, tenant string, item *model.Bundle) error {
	ret := _m.Called(ctx, tenant, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Bundle) error); ok {
		r0 = rf(ctx, tenant, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, tenant, id
func (_m *BundleRepository) Delete(ctx context.Context, tenant string, id string) error {
	ret := _m.Called(ctx, tenant, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: ctx, tenant, id
func (_m *BundleRepository) Exists(ctx context.Context, tenant string, id string) (bool, error) {
	ret := _m.Called(ctx, tenant, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string, string) bool); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, tenant, id
func (_m *BundleRepository) GetByID(ctx context.Context, tenant string, id string) (*model.Bundle, error) {
	ret := _m.Called(ctx, tenant, id)

	var r0 *model.Bundle
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Bundle); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Bundle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenant, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForApplication provides a mock function with given fields: ctx, tenant, id, applicationID
func (_m *BundleRepository) GetForApplication(ctx context.Context, tenant string, id string, applicationID string) (*model.Bundle, error) {
	ret := _m.Called(ctx, tenant, id, applicationID)

	var r0 *model.Bundle
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string) *model.Bundle); ok {
		r0 = rf(ctx, tenant, id, applicationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Bundle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string) error); ok {
		r1 = rf(ctx, tenant, id, applicationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByApplicationIDNoPaging provides a mock function with given fields: ctx, tenantID, appID
func (_m *BundleRepository) ListByApplicationIDNoPaging(ctx context.Context, tenantID string, appID string) ([]*model.Bundle, error) {
	ret := _m.Called(ctx, tenantID, appID)

	var r0 []*model.Bundle
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*model.Bundle); ok {
		r0 = rf(ctx, tenantID, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Bundle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, tenantID, appID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByApplicationIDs provides a mock function with given fields: ctx, tenantID, applicationIDs, pageSize, cursor
func (_m *BundleRepository) ListByApplicationIDs(ctx context.Context, tenantID string, applicationIDs []string, pageSize int, cursor string) ([]*model.BundlePage, error) {
	ret := _m.Called(ctx, tenantID, applicationIDs, pageSize, cursor)

	var r0 []*model.BundlePage
	if rf, ok := ret.Get(0).(func(context.Context, string, []string, int, string) []*model.BundlePage); ok {
		r0 = rf(ctx, tenantID, applicationIDs, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.BundlePage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, []string, int, string) error); ok {
		r1 = rf(ctx, tenantID, applicationIDs, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, tenant, item
func (_m *BundleRepository) Update(ctx context.Context, tenant string, item *model.Bundle) error {
	ret := _m.Called(ctx, tenant, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Bundle) error); ok {
		r0 = rf(ctx, tenant, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewBundleRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewBundleRepository creates a new instance of BundleRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBundleRepository(t mockConstructorTestingTNewBundleRepository) *BundleRepository {
	mock := &BundleRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
