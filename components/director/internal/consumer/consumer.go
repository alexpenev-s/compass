package consumer

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/oathkeeper"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
)

// ConsumerType missing godoc
type ConsumerType string

const (
	// Runtime missing godoc
	Runtime ConsumerType = "Runtime"
	// Application missing godoc
	Application ConsumerType = "Application"
	// IntegrationSystem missing godoc
	IntegrationSystem ConsumerType = "Integration System"
	// User missing godoc
	User ConsumerType = "Static User"
)

// Consumer missing godoc
type Consumer struct {
	ConsumerID string
	ConsumerType
	Flow oathkeeper.AuthFlow
}

// MapSystemAuthToConsumerType missing godoc
func MapSystemAuthToConsumerType(refObj model.SystemAuthReferenceObjectType) (ConsumerType, error) {
	switch refObj {
	case model.ApplicationReference:
		return Application, nil
	case model.RuntimeReference:
		return Runtime, nil
	case model.IntegrationSystemReference:
		return IntegrationSystem, nil
	}
	return "", apperrors.NewInternalError("unknown reference object type")
}
