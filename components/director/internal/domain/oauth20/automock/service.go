// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"
	"github.com/kyma-incubator/compass/components/director/pkg/systemauth"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CreateClientCredentials provides a mock function with given fields: ctx, objectType
func (_m *Service) CreateClientCredentials(ctx context.Context, objectType systemauth.SystemAuthReferenceObjectType) (*model.OAuthCredentialDataInput, error) {
	ret := _m.Called(ctx, objectType)

	var r0 *model.OAuthCredentialDataInput
	if rf, ok := ret.Get(0).(func(context.Context, systemauth.SystemAuthReferenceObjectType) *model.OAuthCredentialDataInput); ok {
		r0 = rf(ctx, objectType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OAuthCredentialDataInput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, systemauth.SystemAuthReferenceObjectType) error); ok {
		r1 = rf(ctx, objectType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteClientCredentials provides a mock function with given fields: ctx, clientID
func (_m *Service) DeleteClientCredentials(ctx context.Context, clientID string) error {
	ret := _m.Called(ctx, clientID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, clientID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
