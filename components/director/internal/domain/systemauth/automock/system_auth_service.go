// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	internalmodel "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/pkg/model"
)

// SystemAuthService is an autogenerated mock type for the SystemAuthService type
type SystemAuthService struct {
	mock.Mock
}

// DeleteByIDForObject provides a mock function with given fields: ctx, objectType, authID
func (_m *SystemAuthService) DeleteByIDForObject(ctx context.Context, objectType model.SystemAuthReferenceObjectType, authID string) error {
	ret := _m.Called(ctx, objectType, authID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.SystemAuthReferenceObjectType, string) error); ok {
		r0 = rf(ctx, objectType, authID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByIDForObject provides a mock function with given fields: ctx, objectType, authID
func (_m *SystemAuthService) GetByIDForObject(ctx context.Context, objectType model.SystemAuthReferenceObjectType, authID string) (*model.SystemAuth, error) {
	ret := _m.Called(ctx, objectType, authID)

	var r0 *model.SystemAuth
	if rf, ok := ret.Get(0).(func(context.Context, model.SystemAuthReferenceObjectType, string) *model.SystemAuth); ok {
		r0 = rf(ctx, objectType, authID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SystemAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.SystemAuthReferenceObjectType, string) error); ok {
		r1 = rf(ctx, objectType, authID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByToken provides a mock function with given fields: ctx, token
func (_m *SystemAuthService) GetByToken(ctx context.Context, token string) (*model.SystemAuth, error) {
	ret := _m.Called(ctx, token)

	var r0 *model.SystemAuth
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.SystemAuth); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SystemAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetGlobal provides a mock function with given fields: ctx, id
func (_m *SystemAuthService) GetGlobal(ctx context.Context, id string) (*model.SystemAuth, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.SystemAuth
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.SystemAuth); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SystemAuth)
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

// InvalidateToken provides a mock function with given fields: ctx, id
func (_m *SystemAuthService) InvalidateToken(ctx context.Context, id string) (*model.SystemAuth, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.SystemAuth
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.SystemAuth); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SystemAuth)
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

// Update provides a mock function with given fields: ctx, item
func (_m *SystemAuthService) Update(ctx context.Context, item *model.SystemAuth) error {
	ret := _m.Called(ctx, item)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.SystemAuth) error); ok {
		r0 = rf(ctx, item)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateValue provides a mock function with given fields: ctx, id, item
func (_m *SystemAuthService) UpdateValue(ctx context.Context, id string, item *internalmodel.Auth) (*model.SystemAuth, error) {
	ret := _m.Called(ctx, id, item)

	var r0 *model.SystemAuth
	if rf, ok := ret.Get(0).(func(context.Context, string, *internalmodel.Auth) *model.SystemAuth); ok {
		r0 = rf(ctx, id, item)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.SystemAuth)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *internalmodel.Auth) error); ok {
		r1 = rf(ctx, id, item)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewSystemAuthService interface {
	mock.TestingT
	Cleanup(func())
}

// NewSystemAuthService creates a new instance of SystemAuthService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewSystemAuthService(t mockConstructorTestingTNewSystemAuthService) *SystemAuthService {
	mock := &SystemAuthService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
