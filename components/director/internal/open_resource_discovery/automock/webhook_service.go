// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
	mock "github.com/stretchr/testify/mock"
)

// WebhookService is an autogenerated mock type for the WebhookService type
type WebhookService struct {
	mock.Mock
}

// ListByWebhookType provides a mock function with given fields: ctx, webhookType
func (_m *WebhookService) ListByWebhookType(ctx context.Context, webhookType model.WebhookType) ([]*model.Webhook, error) {
	ret := _m.Called(ctx, webhookType)

	var r0 []*model.Webhook
	if rf, ok := ret.Get(0).(func(context.Context, model.WebhookType) []*model.Webhook); ok {
		r0 = rf(ctx, webhookType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Webhook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.WebhookType) error); ok {
		r1 = rf(ctx, webhookType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewWebhookService interface {
	mock.TestingT
	Cleanup(func())
}

// NewWebhookService creates a new instance of WebhookService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWebhookService(t mockConstructorTestingTNewWebhookService) *WebhookService {
	mock := &WebhookService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
