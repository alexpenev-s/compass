// Code generated by mockery. DO NOT EDIT.

package automock

import (
	context "context"
	json "encoding/json"

	destinationcreator "github.com/kyma-incubator/compass/components/director/internal/domain/destination/destinationcreator"

	mock "github.com/stretchr/testify/mock"

	model "github.com/kyma-incubator/compass/components/director/internal/model"

	operators "github.com/kyma-incubator/compass/components/director/internal/domain/formationconstraint/operators"
)

// DestinationService is an autogenerated mock type for the destinationService type
type DestinationService struct {
	mock.Mock
}

// CreateBasicCredentialDestinations provides a mock function with given fields: ctx, destinationDetails, basicAuthenticationCredentials, formationAssignment, correlationIDs
func (_m *DestinationService) CreateBasicCredentialDestinations(ctx context.Context, destinationDetails operators.Destination, basicAuthenticationCredentials operators.BasicAuthentication, formationAssignment *model.FormationAssignment, correlationIDs []string) (int, error) {
	ret := _m.Called(ctx, destinationDetails, basicAuthenticationCredentials, formationAssignment, correlationIDs)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, operators.Destination, operators.BasicAuthentication, *model.FormationAssignment, []string) int); ok {
		r0 = rf(ctx, destinationDetails, basicAuthenticationCredentials, formationAssignment, correlationIDs)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, operators.Destination, operators.BasicAuthentication, *model.FormationAssignment, []string) error); ok {
		r1 = rf(ctx, destinationDetails, basicAuthenticationCredentials, formationAssignment, correlationIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateCertificateInDestinationService provides a mock function with given fields: ctx, destinationDetails, formationAssignment
func (_m *DestinationService) CreateCertificateInDestinationService(ctx context.Context, destinationDetails operators.Destination, formationAssignment *model.FormationAssignment) (destinationcreator.CertificateResponse, int, error) {
	ret := _m.Called(ctx, destinationDetails, formationAssignment)

	var r0 destinationcreator.CertificateResponse
	if rf, ok := ret.Get(0).(func(context.Context, operators.Destination, *model.FormationAssignment) destinationcreator.CertificateResponse); ok {
		r0 = rf(ctx, destinationDetails, formationAssignment)
	} else {
		r0 = ret.Get(0).(destinationcreator.CertificateResponse)
	}

	var r1 int
	if rf, ok := ret.Get(1).(func(context.Context, operators.Destination, *model.FormationAssignment) int); ok {
		r1 = rf(ctx, destinationDetails, formationAssignment)
	} else {
		r1 = ret.Get(1).(int)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, operators.Destination, *model.FormationAssignment) error); ok {
		r2 = rf(ctx, destinationDetails, formationAssignment)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// CreateDesignTimeDestinations provides a mock function with given fields: ctx, destinationDetails, formationAssignment
func (_m *DestinationService) CreateDesignTimeDestinations(ctx context.Context, destinationDetails operators.Destination, formationAssignment *model.FormationAssignment) (int, error) {
	ret := _m.Called(ctx, destinationDetails, formationAssignment)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, operators.Destination, *model.FormationAssignment) int); ok {
		r0 = rf(ctx, destinationDetails, formationAssignment)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, operators.Destination, *model.FormationAssignment) error); ok {
		r1 = rf(ctx, destinationDetails, formationAssignment)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateSAMLAssertionDestination provides a mock function with given fields: ctx, destinationDetails, samlAssertionAuthCredentials, formationAssignment, correlationIDs
func (_m *DestinationService) CreateSAMLAssertionDestination(ctx context.Context, destinationDetails operators.Destination, samlAssertionAuthCredentials *operators.SAMLAssertionAuthentication, formationAssignment *model.FormationAssignment, correlationIDs []string) (int, error) {
	ret := _m.Called(ctx, destinationDetails, samlAssertionAuthCredentials, formationAssignment, correlationIDs)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, operators.Destination, *operators.SAMLAssertionAuthentication, *model.FormationAssignment, []string) int); ok {
		r0 = rf(ctx, destinationDetails, samlAssertionAuthCredentials, formationAssignment, correlationIDs)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, operators.Destination, *operators.SAMLAssertionAuthentication, *model.FormationAssignment, []string) error); ok {
		r1 = rf(ctx, destinationDetails, samlAssertionAuthCredentials, formationAssignment, correlationIDs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCertificateFromDestinationService provides a mock function with given fields: ctx, certificateName, externalDestSubaccountID, formationAssignment
func (_m *DestinationService) DeleteCertificateFromDestinationService(ctx context.Context, certificateName string, externalDestSubaccountID string, formationAssignment *model.FormationAssignment) error {
	ret := _m.Called(ctx, certificateName, externalDestSubaccountID, formationAssignment)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.FormationAssignment) error); ok {
		r0 = rf(ctx, certificateName, externalDestSubaccountID, formationAssignment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDestinationFromDestinationService provides a mock function with given fields: ctx, destinationName, destinationSubaccount, formationAssignment
func (_m *DestinationService) DeleteDestinationFromDestinationService(ctx context.Context, destinationName string, destinationSubaccount string, formationAssignment *model.FormationAssignment) error {
	ret := _m.Called(ctx, destinationName, destinationSubaccount, formationAssignment)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *model.FormationAssignment) error); ok {
		r0 = rf(ctx, destinationName, destinationSubaccount, formationAssignment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteDestinations provides a mock function with given fields: ctx, formationAssignment
func (_m *DestinationService) DeleteDestinations(ctx context.Context, formationAssignment *model.FormationAssignment) error {
	ret := _m.Called(ctx, formationAssignment)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.FormationAssignment) error); ok {
		r0 = rf(ctx, formationAssignment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EnrichAssignmentConfigWithCertificateData provides a mock function with given fields: assignmentConfig, certData, destinationIndex
func (_m *DestinationService) EnrichAssignmentConfigWithCertificateData(assignmentConfig json.RawMessage, certData destinationcreator.CertificateResponse, destinationIndex int) (json.RawMessage, error) {
	ret := _m.Called(assignmentConfig, certData, destinationIndex)

	var r0 json.RawMessage
	if rf, ok := ret.Get(0).(func(json.RawMessage, destinationcreator.CertificateResponse, int) json.RawMessage); ok {
		r0 = rf(assignmentConfig, certData, destinationIndex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(json.RawMessage)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(json.RawMessage, destinationcreator.CertificateResponse, int) error); ok {
		r1 = rf(assignmentConfig, certData, destinationIndex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewDestinationService interface {
	mock.TestingT
	Cleanup(func())
}

// NewDestinationService creates a new instance of DestinationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDestinationService(t mockConstructorTestingTNewDestinationService) *DestinationService {
	mock := &DestinationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
