package tombstone

// Entity represents a tombstone entity.
type Entity struct {
	ID            string `db:"id"`
	OrdID         string `db:"ord_id"`
	ApplicationID string `db:"app_id"`
	RemovalDate   string `db:"removal_date"`
}

// GetID returns the entity's ID.
func (e *Entity) GetID() string {
	return e.ID
}

// GetParentID returns the entity's parent ID.
func (e *Entity) GetParentID() string {
	return e.ApplicationID
}

// DecorateWithTenantID decorates the entity with the given tenant ID.
func (e *Entity) DecorateWithTenantID(tenant string) interface{} {
	return struct {
		*Entity
		TenantID string `db:"tenant_id"`
	}{
		Entity:   e,
		TenantID: tenant,
	}
}
