package consumer

import (
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/model"
	"github.com/kyma-incubator/compass/components/hydrator/pkg/oathkeeper"
)

// ConsumerType missing godoc
type ConsumerType string

const (
	// Runtime missing godoc
	Runtime ConsumerType = "Runtime"
	// ExternalCertificate missing godoc
	ExternalCertificate ConsumerType = "External Certificate"
	// Application missing godoc
	Application ConsumerType = "Application"
	// IntegrationSystem missing godoc
	IntegrationSystem ConsumerType = "Integration System"
	// User missing godoc
	User ConsumerType = "Static User"
	// SuperAdmin is a consumer type that is used only in our tests
	SuperAdmin ConsumerType = "Super Admin"
	// TechnicalClient is a consumer type that is used by Atom
	TechnicalClient ConsumerType = "Technical Client"
	// BusinessIntegration missing godoc
	BusinessIntegration ConsumerType = "Business Integration"
	// ManagedApplicationProviderOperator is a consumer type that is used by Managed Application Provider operator
	ManagedApplicationProviderOperator ConsumerType = "Managed Application Provider Operator"
	// ManagedApplicationConsumer is a consumer type that is used by Managed Application Provider operator
	// when creating Certificate Subject Mappings
	ManagedApplicationConsumer ConsumerType = "Managed Application Consumer"
	// LandscapeResourceOperator is a consumer type that is used by Landscape Resource operator
	LandscapeResourceOperator ConsumerType = "Landscape Resource Operator"
)

// Consumer missing godoc
type Consumer struct {
	ConsumerID    string `json:"ConsumerID"`
	ConsumerType  `json:"ConsumerType"`
	Flow          oathkeeper.AuthFlow `json:"Flow"`
	OnBehalfOf    string              `json:"onBehalfOf"`
	Region        string              `json:"region"`
	TokenClientID string              `json:"tokenClientID"`
}

// MapSystemAuthToConsumerType missing godoc
func MapSystemAuthToConsumerType(refObj model.SystemAuthReferenceObjectType) (ConsumerType, error) {
	switch refObj {
	case model.ApplicationReference:
		return Application, nil
	case model.ExternalCertificateReference:
		return ExternalCertificate, nil
	case model.RuntimeReference:
		return Runtime, nil
	case model.IntegrationSystemReference:
		return IntegrationSystem, nil
	case model.TechnicalClientReference:
		return TechnicalClient, nil
	case model.BusinessIntegrationReference:
		return BusinessIntegration, nil
	case model.ManagedApplicationProviderOperatorReference:
		return ManagedApplicationProviderOperator, nil
	case model.ManagedApplicationConsumerReference:
		return ManagedApplicationConsumer, nil
	case model.LandscapeResourceOperatorConsumerReference:
		return LandscapeResourceOperator, nil
	case model.SuperAdminReference:
		return SuperAdmin, nil
	}
	return "", apperrors.NewInternalError("unknown reference object type")
}
