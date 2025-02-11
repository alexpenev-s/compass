package operation_test

import (
	"database/sql"
	"testing"

	"github.com/kyma-incubator/compass/components/director/internal/domain/operation"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntityConverter_ToEntity(t *testing.T) {
	t.Run("success all nullable properties filled", func(t *testing.T) {
		// GIVEN
		opModel := fixOperationModel(ordOpType, model.OperationStatusScheduled)
		require.NotNil(t, opModel)

		conv := operation.NewConverter()

		// WHEN
		entity := conv.ToEntity(opModel)

		// THEN
		expectedOperation := fixEntityOperation(operationID, ordOpType, model.OperationStatusScheduled)

		assert.Equal(t, expectedOperation, entity)
	})
	t.Run("success all nullable properties empty", func(t *testing.T) {
		// GIVEN
		opModel := &model.Operation{
			ID:         operationID,
			OpType:     ordOpType,
			Status:     model.OperationStatusScheduled,
			Data:       nil,
			Error:      nil,
			Priority:   1,
			CreatedAt:  nil,
			FinishedAt: nil,
		}

		expectedEntity := &operation.Entity{
			ID:         operationID,
			Type:       ordOpType,
			Status:     string(model.OperationStatusScheduled),
			Data:       sql.NullString{},
			Error:      sql.NullString{},
			Priority:   1,
			CreatedAt:  nil,
			FinishedAt: nil,
		}
		conv := operation.NewConverter()

		// WHEN
		entity := conv.ToEntity(opModel)

		// THEN
		assert.Equal(t, expectedEntity, entity)
	})
}

func TestEntityConverter_FromEntity(t *testing.T) {
	t.Run("success all nullable properties filled", func(t *testing.T) {
		// GIVEN
		entity := fixEntityOperation(operationID, ordOpType, model.OperationStatusScheduled)
		conv := operation.NewConverter()

		// WHEN
		opModel := conv.FromEntity(entity)

		// THEN
		expectedOperation := fixOperationModel(ordOpType, model.OperationStatusScheduled)
		assert.Equal(t, expectedOperation, opModel)
	})

	t.Run("success all nullable properties empty", func(t *testing.T) {
		// GIVEN
		entity := &operation.Entity{
			ID:         operationID,
			Type:       ordOpType,
			Status:     string(model.OperationStatusScheduled),
			Data:       sql.NullString{},
			Error:      sql.NullString{},
			Priority:   1,
			CreatedAt:  nil,
			FinishedAt: nil,
		}
		expectedModel := &model.Operation{
			ID:         operationID,
			OpType:     ordOpType,
			Status:     model.OperationStatusScheduled,
			Data:       nil,
			Error:      nil,
			Priority:   1,
			CreatedAt:  nil,
			FinishedAt: nil,
		}
		conv := operation.NewConverter()

		// WHEN
		opModel := conv.FromEntity(entity)

		// THEN
		assert.Equal(t, expectedModel, opModel)
	})
}
