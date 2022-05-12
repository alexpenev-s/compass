// Code generated by mockery v2.12.1. DO NOT EDIT.

package automock

import (
	testing "testing"

	mock "github.com/stretchr/testify/mock"

	tls "crypto/tls"
)

// CertificateCache is an autogenerated mock type for the CertificateCache type
type CertificateCache struct {
	mock.Mock
}

// Get provides a mock function with given fields:
func (_m *CertificateCache) Get() *tls.Certificate {
	ret := _m.Called()

	var r0 *tls.Certificate
	if rf, ok := ret.Get(0).(func() *tls.Certificate); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tls.Certificate)
		}
	}

	return r0
}

// NewCertificateCache creates a new instance of CertificateCache. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewCertificateCache(t testing.TB) *CertificateCache {
	mock := &CertificateCache{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
