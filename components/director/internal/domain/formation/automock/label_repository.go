// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"
)

// LabelRepository is an autogenerated mock type for the labelRepository type
type LabelRepository struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4
func (_m *LabelRepository) Delete(_a0 context.Context, _a1 string, _a2 model.LabelableObject, _a3 string, _a4 string) error {
	ret := _m.Called(_a0, _a1, _a2, _a3, _a4)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, model.LabelableObject, string, string) error); ok {
		r0 = rf(_a0, _a1, _a2, _a3, _a4)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type NewLabelRepositoryT interface {
	mock.TestingT
	Cleanup(func())
}

// NewLabelRepository creates a new instance of LabelRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLabelRepository(t NewLabelRepositoryT) *LabelRepository {
	mock := &LabelRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
