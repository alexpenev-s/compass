package operation

import (
	"database/sql"
	"time"
)

// Entity is a representation of an Operation in the database.
type Entity struct {
	ID         string         `db:"id"`
	Type       string         `db:"op_type"`
	Status     string         `db:"status"`
	Data       sql.NullString `db:"data"`
	Error      sql.NullString `db:"error"`
	Priority   int            `db:"priority"`
	CreatedAt  *time.Time     `db:"created_at"`
	FinishedAt *time.Time     `db:"finished_at"`
}
