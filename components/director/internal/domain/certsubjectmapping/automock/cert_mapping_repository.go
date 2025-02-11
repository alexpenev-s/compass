// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// CertMappingRepository is an autogenerated mock type for the CertMappingRepository type
type CertMappingRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, item
func (_m *CertMappingRepository) Create(ctx context.Context, item *model.CertSubjectMapping) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.CertSubjectMapping) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *CertMappingRepository) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteByConsumerID provides a mock function with given fields: ctx, consumerID
func (_m *CertMappingRepository) DeleteByConsumerID(ctx context.Context, consumerID string) error {
	ret := _m.Called(ctx, consumerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, consumerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exists provides a mock function with given fields: ctx, id
func (_m *CertMappingRepository) Exists(ctx context.Context, id string) (bool, error) {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, id
func (_m *CertMappingRepository) Get(ctx context.Context, id string) (*model.CertSubjectMapping, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.CertSubjectMapping
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.CertSubjectMapping); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.CertSubjectMapping)
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

// List provides a mock function with given fields: ctx, pageSize, cursor
func (_m *CertMappingRepository) List(ctx context.Context, pageSize int, cursor string) (*model.CertSubjectMappingPage, error) {
	ret := _m.Called(ctx, pageSize, cursor)

	var r0 *model.CertSubjectMappingPage
	if rf, ok := ret.Get(0).(func(context.Context, int, string) *model.CertSubjectMappingPage); ok {
		r0 = rf(ctx, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.CertSubjectMappingPage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, string) error); ok {
		r1 = rf(ctx, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, _a1
func (_m *CertMappingRepository) Update(ctx context.Context, _a1 *model.CertSubjectMapping) error {
	ret := _m.Called(ctx, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.CertSubjectMapping) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCertMappingRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewCertMappingRepository creates a new instance of CertMappingRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCertMappingRepository(t mockConstructorTestingTNewCertMappingRepository) *CertMappingRepository {
	mock := &CertMappingRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
