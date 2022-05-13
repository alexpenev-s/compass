// Code generated by mockery v2.12.2. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// PackageRepository is an autogenerated mock type for the PackageRepository type
type PackageRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, tenant, item
func (_m *PackageRepository) Create(ctx context.Context, tenant string, item *model.Package) error {
	ret := _m.Called(ctx, tenant, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Package) error); ok {
		r0 = rf(ctx, tenant, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, tenant, id
func (_m *PackageRepository) Delete(ctx context.Context, tenant string, id string) error {
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
func (_m *PackageRepository) Exists(ctx context.Context, tenant string, id string) (bool, error) {
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
func (_m *PackageRepository) GetByID(ctx context.Context, tenant string, id string) (*model.Package, error) {
	ret := _m.Called(ctx, tenant, id)

	var r0 *model.Package
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Package); ok {
		r0 = rf(ctx, tenant, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Package)
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

// ListByApplicationID provides a mock function with given fields: ctx, tenantID, appID
func (_m *PackageRepository) ListByApplicationID(ctx context.Context, tenantID string, appID string) ([]*model.Package, error) {
	ret := _m.Called(ctx, tenantID, appID)

	var r0 []*model.Package
	if rf, ok := ret.Get(0).(func(context.Context, string, string) []*model.Package); ok {
		r0 = rf(ctx, tenantID, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Package)
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

// Update provides a mock function with given fields: ctx, tenant, item
func (_m *PackageRepository) Update(ctx context.Context, tenant string, item *model.Package) error {
	ret := _m.Called(ctx, tenant, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.Package) error); ok {
		r0 = rf(ctx, tenant, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPackageRepository creates a new instance of PackageRepository. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewPackageRepository(t testing.TB) *PackageRepository {
	mock := &PackageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
