// Code generated by mockery v1.0.0. DO NOT EDIT.

package automock

import (
	context "context"

	model "github.com/kyma-incubator/compass/components/director/internal2/model"
	mock "github.com/stretchr/testify/mock"
)

// LabelUpsertService is an autogenerated mock type for the LabelUpsertService type
type LabelUpsertService struct {
	mock.Mock
}

// UpsertLabel provides a mock function with given fields: ctx, tenant, labelInput
func (_m *LabelUpsertService) UpsertLabel(ctx context.Context, tenant string, labelInput *model.LabelInput) error {
	ret := _m.Called(ctx, tenant, labelInput)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *model.LabelInput) error); ok {
		r0 = rf(ctx, tenant, labelInput)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
