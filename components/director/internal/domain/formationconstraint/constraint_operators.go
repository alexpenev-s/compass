package formationconstraint

import (
	"context"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/internal/domain/label"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/formationconstraint"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"github.com/pkg/errors"
)

const (
	// IsNotAssignedToAnyFormationOfTypeOperator represents the IsNotAssignedToAnyFormationOfType operator
	IsNotAssignedToAnyFormationOfTypeOperator = "IsNotAssignedToAnyFormationOfType"
	// DoesNotContainResourceOfSubtypeOperator represents the DoesNotContainResourceOfSubtype operator
	DoesNotContainResourceOfSubtypeOperator = "DoesNotContainResourceOfSubtype"
)

// OperatorName represents the constraint operator name
type OperatorName string

// OperatorInput represents the input needed by the constraint operator
type OperatorInput interface{}

// OperatorFunc provides an interface for functions implementing constraint operators
type OperatorFunc func(ctx context.Context, input OperatorInput) (bool, error)

// OperatorInputConstructor returns empty OperatorInput for a certain constraint operator
type OperatorInputConstructor func() OperatorInput

// NewIsNotAssignedToAnyFormationOfTypeInput is input constructor for IsNotAssignedToAnyFormationOfType operator. It returns empty OperatorInput.
func NewIsNotAssignedToAnyFormationOfTypeInput() OperatorInput {
	return &formationconstraint.IsNotAssignedToAnyFormationOfTypeInput{}
}

// NewDoesNotContainResourceOfSubtypeInput is input constructor for DoesNotContainResourceOfSubtypeOperator operator. It returns empty OperatorInput
func NewDoesNotContainResourceOfSubtypeInput() OperatorInput {
	return &formationconstraint.DoesNotContainResourceOfSubtypeInput{}
}

// IsNotAssignedToAnyFormationOfType is a constraint operator. It checks if the resource from the OperatorInput is already part of formation of the type that the operator is associated with
func (e *ConstraintEngine) IsNotAssignedToAnyFormationOfType(ctx context.Context, input OperatorInput) (bool, error) {
	log.C(ctx).Infof("Executing operator: %s", IsNotAssignedToAnyFormationOfTypeOperator)

	i, ok := input.(*formationconstraint.IsNotAssignedToAnyFormationOfTypeInput)
	if !ok {
		return false, errors.New("Incompatible input")
	}

	log.C(ctx).Infof("Enforcing constraint on resource of type: %q, subtype: %q and ID: %q", i.ResourceType, i.ResourceSubtype, i.ResourceID)

	var assignedFormations []string
	switch i.ResourceType {
	case model.TenantResourceType:
		tenantInternalID, err := e.tenantSvc.GetInternalTenant(ctx, i.ResourceID)
		if err != nil {
			return false, err
		}

		assignments, err := e.asaSvc.ListForTargetTenant(ctx, tenantInternalID)
		if err != nil {
			return false, err
		}

		assignedFormations = make([]string, 0, len(assignments))
		for _, a := range assignments {
			assignedFormations = append(assignedFormations, a.ScenarioName)
		}
	case model.ApplicationResourceType:
		scenariosLabel, err := e.labelRepo.GetByKey(ctx, i.Tenant, model.ApplicationLabelableObject, i.ResourceID, model.ScenariosKey)
		if err != nil {
			if apperrors.IsNotFoundError(err) {
				return true, nil
			}
			return false, err
		}
		assignedFormations, err = label.ValueToStringsSlice(scenariosLabel.Value)
		if err != nil {
			return false, err
		}
	default:
		return false, errors.Errorf("Unsupported resource type %q", i.ResourceType)
	}

	isAllowedToParticipateInFormationsOfType, err := e.isAllowedToParticipateInFormationsOfType(ctx, assignedFormations, i.Tenant, i.FormationTemplateID, i.ResourceSubtype, i.ExceptSystemTypes)
	if err != nil {
		return false, err
	}

	return isAllowedToParticipateInFormationsOfType, nil
}

// DoesNotContainResourceOfSubtype is a constraint operator. It checks if the formation contains resource with the same subtype as the resource subtype from the OperatorInput
func (e *ConstraintEngine) DoesNotContainResourceOfSubtype(ctx context.Context, input OperatorInput) (bool, error) {
	log.C(ctx).Infof("Executing operator: %q", DoesNotContainResourceOfSubtypeOperator)

	i, ok := input.(*formationconstraint.DoesNotContainResourceOfSubtypeInput)
	if !ok {
		return false, errors.New(fmt.Sprintf("Incompatible input for operator %q", DoesNotContainResourceOfSubtypeOperator))
	}

	log.C(ctx).Infof("Enforcing constraint on resource of type: %q, subtype: %q and ID: %q", i.ResourceType, i.ResourceSubtype, i.ResourceID)

	switch i.ResourceType {
	case model.ApplicationResourceType:
		applications, err := e.applicationRepository.ListByScenariosNoPaging(ctx, i.Tenant, []string{i.FormationName})
		if err != nil {
			return false, errors.Wrap(err, fmt.Sprintf("while listing applications in scenario %q", i.FormationName))
		}

		for _, application := range applications {
			appTypeLbl, err := e.labelService.GetByKey(ctx, i.Tenant, model.ApplicationLabelableObject, application.ID, i.ResourceTypeLabelKey)
			if err != nil {
				return false, errors.Wrap(err, fmt.Sprintf("while getting label with key %q of application with ID %q in tenant %q", i.ResourceTypeLabelKey, application.ID, i.Tenant))
			}

			if i.ResourceSubtype == appTypeLbl.Value.(string) {
				return false, nil
			}
		}
	default:
		return false, errors.Errorf("Unsupported resource type %q", i.ResourceType)
	}

	return true, nil
}

func (e *ConstraintEngine) isAllowedToParticipateInFormationsOfType(ctx context.Context, assignedFormationNames []string, tenant, formationTemplateID, resourceSubtype string, exceptSystemTypes []string) (bool, error) {
	if len(assignedFormationNames) == 0 {
		return true, nil
	}

	if len(exceptSystemTypes) > 0 {
		for _, exceptType := range exceptSystemTypes {
			if resourceSubtype == exceptType {
				return true, nil
			}
		}
	}

	assignedFormations, err := e.formationRepo.ListByFormationNames(ctx, assignedFormationNames, tenant)
	if err != nil {
		return false, err
	}

	for _, formation := range assignedFormations {
		if formation.FormationTemplateID == formationTemplateID {
			return false, nil
		}
	}

	return true, nil
}
