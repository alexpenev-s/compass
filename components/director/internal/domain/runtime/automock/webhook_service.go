// Code generated by mockery v2.12.2. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// WebhookService is an autogenerated mock type for the WebhookService type
type WebhookService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, owningResourceID, in, objectType
func (_m *WebhookService) Create(ctx context.Context, owningResourceID string, in model.WebhookInput, objectType model.WebhookReferenceObjectType) (string, error) {
	ret := _m.Called(ctx, owningResourceID, in, objectType)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string, model.WebhookInput, model.WebhookReferenceObjectType) string); ok {
		r0 = rf(ctx, owningResourceID, in, objectType)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, model.WebhookInput, model.WebhookReferenceObjectType) error); ok {
		r1 = rf(ctx, owningResourceID, in, objectType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListForRuntime provides a mock function with given fields: ctx, runtimeID
func (_m *WebhookService) ListForRuntime(ctx context.Context, runtimeID string) ([]*model.Webhook, error) {
	ret := _m.Called(ctx, runtimeID)

	var r0 []*model.Webhook
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Webhook); ok {
		r0 = rf(ctx, runtimeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Webhook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, runtimeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewWebhookService creates a new instance of WebhookService. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewWebhookService(t testing.TB) *WebhookService {
	mock := &WebhookService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
