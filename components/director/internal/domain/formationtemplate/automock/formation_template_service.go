// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// FormationTemplateService is an autogenerated mock type for the FormationTemplateService type
type FormationTemplateService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in
func (_m *FormationTemplateService) Create(ctx context.Context, in *model.FormationTemplateInput) (string, error) {
	ret := _m.Called(ctx, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, *model.FormationTemplateInput) string); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *model.FormationTemplateInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *FormationTemplateService) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, id
func (_m *FormationTemplateService) Get(ctx context.Context, id string) (*model.FormationTemplate, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.FormationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.FormationTemplate); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationTemplate)
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
func (_m *FormationTemplateService) List(ctx context.Context, pageSize int, cursor string) (*model.FormationTemplatePage, error) {
	ret := _m.Called(ctx, pageSize, cursor)

	var r0 *model.FormationTemplatePage
	if rf, ok := ret.Get(0).(func(context.Context, int, string) *model.FormationTemplatePage); ok {
		r0 = rf(ctx, pageSize, cursor)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FormationTemplatePage)
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

// ListWebhooksForFormationTemplate provides a mock function with given fields: ctx, formationTemplateID
func (_m *FormationTemplateService) ListWebhooksForFormationTemplate(ctx context.Context, formationTemplateID string) ([]*model.Webhook, error) {
	ret := _m.Called(ctx, formationTemplateID)

	var r0 []*model.Webhook
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.Webhook); ok {
		r0 = rf(ctx, formationTemplateID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Webhook)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, formationTemplateID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, in
func (_m *FormationTemplateService) Update(ctx context.Context, id string, in *model.FormationTemplateInput) error {
	ret := _m.Called(ctx, id, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.FormationTemplateInput) error); ok {
		r0 = rf(ctx, id, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewFormationTemplateService interface {
	mock.TestingT
	Cleanup(func())
}

// NewFormationTemplateService creates a new instance of FormationTemplateService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFormationTemplateService(t mockConstructorTestingTNewFormationTemplateService) *FormationTemplateService {
	mock := &FormationTemplateService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
