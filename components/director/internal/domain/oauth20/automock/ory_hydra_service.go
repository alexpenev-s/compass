// Code generated by mockery. DO NOT EDIT.

package automock

import (
	admin "github.com/ory/hydra-client-go/client/admin"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// OryHydraService is an autogenerated mock type for the OryHydraService type
type OryHydraService struct {
	mock.Mock
}

// CreateOAuth2Client provides a mock function with given fields: params, opts
func (_m *OryHydraService) CreateOAuth2Client(params *admin.CreateOAuth2ClientParams, opts ...admin.ClientOption) (*admin.CreateOAuth2ClientCreated, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *admin.CreateOAuth2ClientCreated
	if rf, ok := ret.Get(0).(func(*admin.CreateOAuth2ClientParams, ...admin.ClientOption) *admin.CreateOAuth2ClientCreated); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*admin.CreateOAuth2ClientCreated)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*admin.CreateOAuth2ClientParams, ...admin.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOAuth2Client provides a mock function with given fields: params, opts
func (_m *OryHydraService) DeleteOAuth2Client(params *admin.DeleteOAuth2ClientParams, opts ...admin.ClientOption) (*admin.DeleteOAuth2ClientNoContent, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *admin.DeleteOAuth2ClientNoContent
	if rf, ok := ret.Get(0).(func(*admin.DeleteOAuth2ClientParams, ...admin.ClientOption) *admin.DeleteOAuth2ClientNoContent); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*admin.DeleteOAuth2ClientNoContent)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*admin.DeleteOAuth2ClientParams, ...admin.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListOAuth2Clients provides a mock function with given fields: params, opts
func (_m *OryHydraService) ListOAuth2Clients(params *admin.ListOAuth2ClientsParams, opts ...admin.ClientOption) (*admin.ListOAuth2ClientsOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *admin.ListOAuth2ClientsOK
	if rf, ok := ret.Get(0).(func(*admin.ListOAuth2ClientsParams, ...admin.ClientOption) *admin.ListOAuth2ClientsOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*admin.ListOAuth2ClientsOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*admin.ListOAuth2ClientsParams, ...admin.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOAuth2Client provides a mock function with given fields: params, opts
func (_m *OryHydraService) UpdateOAuth2Client(params *admin.UpdateOAuth2ClientParams, opts ...admin.ClientOption) (*admin.UpdateOAuth2ClientOK, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *admin.UpdateOAuth2ClientOK
	if rf, ok := ret.Get(0).(func(*admin.UpdateOAuth2ClientParams, ...admin.ClientOption) *admin.UpdateOAuth2ClientOK); ok {
		r0 = rf(params, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*admin.UpdateOAuth2ClientOK)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*admin.UpdateOAuth2ClientParams, ...admin.ClientOption) error); ok {
		r1 = rf(params, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOryHydraService creates a new instance of OryHydraService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewOryHydraService(t testing.TB) *OryHydraService {
	mock := &OryHydraService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
