package destination

import (
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/kyma-incubator/compass/components/director/internal/repo"
)

// NewConverter creates a new destination converter
func NewConverter() *converter {
	return &converter{}
}

type converter struct{}

// ToEntity converts from internal model to entity
func (c *converter) ToEntity(in *model.Destination) *Entity {
	if in == nil {
		return nil
	}

	return &Entity{
		ID:                    in.ID,
		Name:                  in.Name,
		Type:                  in.Type,
		URL:                   in.URL,
		Authentication:        in.Authentication,
		TenantID:              in.SubaccountID,
		FormationAssignmentID: repo.NewNullableString(in.FormationAssignmentID),
	}
}

// FromEntity converts from entity to internal model
func (c *converter) FromEntity(e *Entity) *model.Destination {
	if e == nil {
		return nil
	}

	return &model.Destination{
		ID:                    e.ID,
		Name:                  e.Name,
		Type:                  e.Type,
		URL:                   e.URL,
		Authentication:        e.Authentication,
		SubaccountID:          e.TenantID,
		FormationAssignmentID: repo.StringPtrFromNullableString(e.FormationAssignmentID),
	}
}
