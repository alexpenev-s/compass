// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	resource "github.com/kyma-incubator/compass/components/director/pkg/resource"
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

// CreateGlobal provides a mock function with given fields: ctx, _a1
func (_m *PackageRepository) CreateGlobal(ctx context.Context, _a1 *model.Package) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Package) error); ok {
		r0 = rf(ctx, _a1)
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

// DeleteGlobal provides a mock function with given fields: ctx, id
func (_m *PackageRepository) DeleteGlobal(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
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

// GetByIDGlobal provides a mock function with given fields: ctx, id
func (_m *PackageRepository) GetByIDGlobal(ctx context.Context, id string) (*model.Package, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.Package
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.Package); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Package)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByResourceID provides a mock function with given fields: ctx, tenantID, resourceID, resourceType
func (_m *PackageRepository) ListByResourceID(ctx context.Context, tenantID string, resourceID string, resourceType resource.Type) ([]*model.Package, error) {
	ret := _m.Called(ctx, tenantID, resourceID, resourceType)

	var r0 []*model.Package
	if rf, ok := ret.Get(0).(func(context.Context, string, string, resource.Type) []*model.Package); ok {
		r0 = rf(ctx, tenantID, resourceID, resourceType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Package)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, resource.Type) error); ok {
		r1 = rf(ctx, tenantID, resourceID, resourceType)
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

// UpdateGlobal provides a mock function with given fields: ctx, _a1
func (_m *PackageRepository) UpdateGlobal(ctx context.Context, _a1 *model.Package) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Package) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewPackageRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewPackageRepository creates a new instance of PackageRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPackageRepository(t mockConstructorTestingTNewPackageRepository) *PackageRepository {
	mock := &PackageRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
