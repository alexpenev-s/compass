// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	labelfilter "github.com/kyma-incubator/compass/components/director/internal/labelfilter"
	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// ApplicationTemplateService is an autogenerated mock type for the ApplicationTemplateService type
type ApplicationTemplateService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, in
func (_m *ApplicationTemplateService) Create(ctx context.Context, in model.ApplicationTemplateInput) (string, error) {
	ret := _m.Called(ctx, in)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, model.ApplicationTemplateInput) string); ok {
		r0 = rf(ctx, in)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.ApplicationTemplateInput) error); ok {
		r1 = rf(ctx, in)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateWithLabels provides a mock function with given fields: ctx, in, labels
func (_m *ApplicationTemplateService) CreateWithLabels(ctx context.Context, in model.ApplicationTemplateInput, labels map[string]interface{}) (string, error) {
	ret := _m.Called(ctx, in, labels)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, model.ApplicationTemplateInput, map[string]interface{}) string); ok {
		r0 = rf(ctx, in, labels)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, model.ApplicationTemplateInput, map[string]interface{}) error); ok {
		r1 = rf(ctx, in, labels)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *ApplicationTemplateService) Delete(ctx context.Context, id string) error {
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
func (_m *ApplicationTemplateService) Get(ctx context.Context, id string) (*model.ApplicationTemplate, error) {
	ret := _m.Called(ctx, id)

	var r0 *model.ApplicationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, string) *model.ApplicationTemplate); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationTemplate)
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

// GetByFilters provides a mock function with given fields: ctx, filter
func (_m *ApplicationTemplateService) GetByFilters(ctx context.Context, filter []*labelfilter.LabelFilter) (*model.ApplicationTemplate, error) {
	ret := _m.Called(ctx, filter)

	var r0 *model.ApplicationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, []*labelfilter.LabelFilter) *model.ApplicationTemplate); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*labelfilter.LabelFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByNameAndRegion provides a mock function with given fields: ctx, name, region
func (_m *ApplicationTemplateService) GetByNameAndRegion(ctx context.Context, name string, region interface{}) (*model.ApplicationTemplate, error) {
	ret := _m.Called(ctx, name, region)

	var r0 *model.ApplicationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) *model.ApplicationTemplate); ok {
		r0 = rf(ctx, name, region)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.ApplicationTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, interface{}) error); ok {
		r1 = rf(ctx, name, region)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLabel provides a mock function with given fields: ctx, appTemplateID, key
func (_m *ApplicationTemplateService) GetLabel(ctx context.Context, appTemplateID string, key string) (*model.Label, error) {
	ret := _m.Called(ctx, appTemplateID, key)

	var r0 *model.Label
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *model.Label); ok {
		r0 = rf(ctx, appTemplateID, key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, appTemplateID, key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: ctx, filter, pageSize, cursor
func (_m *ApplicationTemplateService) List(ctx context.Context, filter []*labelfilter.LabelFilter, pageSize int, cursor string) (model.ApplicationTemplatePage, error) {
	ret := _m.Called(ctx, filter, pageSize, cursor)

	var r0 model.ApplicationTemplatePage
	if rf, ok := ret.Get(0).(func(context.Context, []*labelfilter.LabelFilter, int, string) model.ApplicationTemplatePage); ok {
		r0 = rf(ctx, filter, pageSize, cursor)
	} else {
		r0 = ret.Get(0).(model.ApplicationTemplatePage)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*labelfilter.LabelFilter, int, string) error); ok {
		r1 = rf(ctx, filter, pageSize, cursor)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByFilters provides a mock function with given fields: ctx, filter
func (_m *ApplicationTemplateService) ListByFilters(ctx context.Context, filter []*labelfilter.LabelFilter) ([]*model.ApplicationTemplate, error) {
	ret := _m.Called(ctx, filter)

	var r0 []*model.ApplicationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, []*labelfilter.LabelFilter) []*model.ApplicationTemplate); ok {
		r0 = rf(ctx, filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ApplicationTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []*labelfilter.LabelFilter) error); ok {
		r1 = rf(ctx, filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListByName provides a mock function with given fields: ctx, name
func (_m *ApplicationTemplateService) ListByName(ctx context.Context, name string) ([]*model.ApplicationTemplate, error) {
	ret := _m.Called(ctx, name)

	var r0 []*model.ApplicationTemplate
	if rf, ok := ret.Get(0).(func(context.Context, string) []*model.ApplicationTemplate); ok {
		r0 = rf(ctx, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ApplicationTemplate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListLabels provides a mock function with given fields: ctx, appTemplateID
func (_m *ApplicationTemplateService) ListLabels(ctx context.Context, appTemplateID string) (map[string]*model.Label, error) {
	ret := _m.Called(ctx, appTemplateID)

	var r0 map[string]*model.Label
	if rf, ok := ret.Get(0).(func(context.Context, string) map[string]*model.Label); ok {
		r0 = rf(ctx, appTemplateID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]*model.Label)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, appTemplateID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PrepareApplicationCreateInputJSON provides a mock function with given fields: appTemplate, values
func (_m *ApplicationTemplateService) PrepareApplicationCreateInputJSON(appTemplate *model.ApplicationTemplate, values model.ApplicationFromTemplateInputValues) (string, error) {
	ret := _m.Called(appTemplate, values)

	var r0 string
	if rf, ok := ret.Get(0).(func(*model.ApplicationTemplate, model.ApplicationFromTemplateInputValues) string); ok {
		r0 = rf(appTemplate, values)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.ApplicationTemplate, model.ApplicationFromTemplateInputValues) error); ok {
		r1 = rf(appTemplate, values)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, in
func (_m *ApplicationTemplateService) Update(ctx context.Context, id string, in model.ApplicationTemplateUpdateInput) error {
	ret := _m.Called(ctx, id, in)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.ApplicationTemplateUpdateInput) error); ok {
		r0 = rf(ctx, id, in)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewApplicationTemplateService interface {
	mock.TestingT
	Cleanup(func())
}

// NewApplicationTemplateService creates a new instance of ApplicationTemplateService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewApplicationTemplateService(t mockConstructorTestingTNewApplicationTemplateService) *ApplicationTemplateService {
	mock := &ApplicationTemplateService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
