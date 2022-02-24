// Code generated by mockery v2.9.4. DO NOT EDIT.

package automock

import (
	context "context"

	authnmappinghandler "github.com/kyma-incubator/compass/components/director/internal/authnmappinghandler"

	mock "github.com/stretchr/testify/mock"
)

// TokenVerifier is an autogenerated mock type for the TokenVerifier type
type TokenVerifier struct {
	mock.Mock
}

// Verify provides a mock function with given fields: ctx, token
func (_m *TokenVerifier) Verify(ctx context.Context, token string) (authnmappinghandler.TokenData, error) {
	ret := _m.Called(ctx, token)

	var r0 authnmappinghandler.TokenData
	if rf, ok := ret.Get(0).(func(context.Context, string) authnmappinghandler.TokenData); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(authnmappinghandler.TokenData)
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
