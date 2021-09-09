package testdb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

// MockDatabase missing godoc
func MockDatabase(t *testing.T) (*sqlx.DB, DBMock) {
	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	sqlxDB := sqlx.NewDb(sqlDB, "sqlmock")

	return sqlxDB, &sqlMockWithAssertions{sqlMock}
}

// DBMock missing godoc
type DBMock interface {
	sqlmock.Sqlmock
	AssertExpectations(t *testing.T)
}

type sqlMockWithAssertions struct {
	sqlmock.Sqlmock
}

// AssertExpectations missing godoc
func (s *sqlMockWithAssertions) AssertExpectations(t *testing.T) {
	err := s.ExpectationsWereMet()
	require.NoError(t, err)
}
