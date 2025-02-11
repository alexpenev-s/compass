package systemssync_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/kyma-incubator/compass/components/director/internal/domain/systemssync"
	"github.com/kyma-incubator/compass/components/director/internal/domain/systemssync/automock"
	"github.com/kyma-incubator/compass/components/director/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_List(t *testing.T) {
	ctx := context.TODO()

	syncTimestamps := []*model.SystemSynchronizationTimestamp{
		fixSystemsSyncModel("id1", "tenant1", "pr", time.Now()),
		fixSystemsSyncModel("id2", "tenant2", "PR", time.Now()),
	}

	testCases := []struct {
		Name           string
		RepositoryFn   func() *automock.SystemsSyncRepository
		ExpectedResult []*model.SystemSynchronizationTimestamp
		ExpectedErr    error
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.SystemsSyncRepository {
				repo := &automock.SystemsSyncRepository{}
				repo.On("List", ctx).Return(syncTimestamps, nil).Once()
				return repo
			},
			ExpectedResult: syncTimestamps,
			ExpectedErr:    nil,
		},
		{
			Name: "Error when listing systems sync timestamps",
			RepositoryFn: func() *automock.SystemsSyncRepository {
				repo := &automock.SystemsSyncRepository{}
				repo.On("List", ctx).Return(nil, testError).Once()
				return repo
			},
			ExpectedResult: nil,
			ExpectedErr:    testError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()

			svc := systemssync.NewService(repo)

			// WHEN
			result, err := svc.List(ctx)

			// THEN
			if testCase.ExpectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.ExpectedResult, result)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestService_Upsert(t *testing.T) {
	ctx := context.TODO()
	syncTimestamp := fixSystemsSyncModel("id1", "tenant1", "pr", time.Now())

	testCases := []struct {
		Name           string
		RepositoryFn   func() *automock.SystemsSyncRepository
		GetInput       func() *model.SystemSynchronizationTimestamp
		ExpectedResult []*model.SystemSynchronizationTimestamp
		ExpectedErr    error
	}{
		{
			Name: "Success",
			RepositoryFn: func() *automock.SystemsSyncRepository {
				repo := &automock.SystemsSyncRepository{}
				repo.On("Upsert", ctx, mock.Anything).Return(nil).Once()
				return repo
			},
			GetInput: func() *model.SystemSynchronizationTimestamp {
				return syncTimestamp
			},
			ExpectedResult: nil,
			ExpectedErr:    nil,
		},
		{
			Name: "Error when upserting systems sync timestamps",
			RepositoryFn: func() *automock.SystemsSyncRepository {
				repo := &automock.SystemsSyncRepository{}
				repo.On("Upsert", ctx, mock.Anything).Return(testError).Once()
				return repo
			},
			GetInput: func() *model.SystemSynchronizationTimestamp {
				return syncTimestamp
			},
			ExpectedResult: nil,
			ExpectedErr:    testError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			repo := testCase.RepositoryFn()

			svc := systemssync.NewService(repo)

			// WHEN
			err := svc.Upsert(ctx, testCase.GetInput())

			// THEN
			if testCase.ExpectedErr != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), testCase.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
