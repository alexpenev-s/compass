// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// BundleRepo is an autogenerated mock type for the BundleRepo type
type BundleRepo struct {
	mock.Mock
}

// GetBySystemAndCorrelationID provides a mock function with given fields: ctx, tenantID, systemName, systemURL, correlationID
func (_m *BundleRepo) GetBySystemAndCorrelationID(ctx context.Context, tenantID string, systemName string, systemURL string, correlationID string) ([]*model.Bundle, error) {
	ret := _m.Called(ctx, tenantID, systemName, systemURL, correlationID)

	var r0 []*model.Bundle
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, string) []*model.Bundle); ok {
		r0 = rf(ctx, tenantID, systemName, systemURL, correlationID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Bundle)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, string) error); ok {
		r1 = rf(ctx, tenantID, systemName, systemURL, correlationID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewBundleRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewBundleRepo creates a new instance of BundleRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBundleRepo(t mockConstructorTestingTNewBundleRepo) *BundleRepo {
	mock := &BundleRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
