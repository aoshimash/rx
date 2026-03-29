package memory

import (
	"context"
	"testing"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlanStore_CreateAndGet(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()
	userID := "user-1"

	plan := &domain.Plan{
		Sessions: []domain.PlanSession{
			{
				SessionName: "Day 1",
				Order:       0,
				Entries: []domain.PlanSessionEntry{
					{ExerciseName: "Squat", Order: 0, Fields: map[string]interface{}{"reps": float64(5)}},
				},
			},
		},
	}

	err := store.Create(ctx, userID, plan)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, plan.ID)

	got, err := store.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, plan.ID, got.ID)
	assert.Len(t, got.Sessions, 1)
	assert.Equal(t, "Day 1", got.Sessions[0].SessionName)
	assert.Len(t, got.Sessions[0].Entries, 1)
	assert.Equal(t, "Squat", got.Sessions[0].Entries[0].ExerciseName)
	assert.NotEqual(t, uuid.Nil, got.Sessions[0].ID)
	assert.NotEqual(t, uuid.Nil, got.Sessions[0].Entries[0].ID)
}

func TestPlanStore_CreateConflict(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()
	userID := "user-1"

	plan := &domain.Plan{}
	require.NoError(t, store.Create(ctx, userID, plan))

	err := store.Create(ctx, userID, &domain.Plan{})
	require.Error(t, err)
	domErr, ok := err.(*domain.DomainError)
	require.True(t, ok)
	assert.Equal(t, domain.ErrorCodeConflict, domErr.Code)
}

func TestPlanStore_Delete(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()
	userID := "user-1"

	plan := &domain.Plan{}
	require.NoError(t, store.Create(ctx, userID, plan))

	err := store.Delete(ctx, userID)
	require.NoError(t, err)

	_, err = store.GetByUserID(ctx, userID)
	assert.ErrorIs(t, err, domain.ErrNotFound)

	// Delete again should return not found
	err = store.Delete(ctx, userID)
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPlanStore_GetNotFound(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()

	_, err := store.GetByUserID(ctx, "nonexistent")
	assert.ErrorIs(t, err, domain.ErrNotFound)
}

func TestPlanStore_Update(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()
	userID := "user-1"

	plan := &domain.Plan{
		Sessions: []domain.PlanSession{
			{SessionName: "Day 1", Order: 0},
		},
	}
	require.NoError(t, store.Create(ctx, userID, plan))
	originalID := plan.ID
	originalCreatedAt := plan.CreatedAt

	updated := &domain.Plan{
		Sessions: []domain.PlanSession{
			{SessionName: "Updated Day 1", Order: 0},
			{SessionName: "Day 2", Order: 1},
		},
	}
	err := store.Update(ctx, userID, updated)
	require.NoError(t, err)

	got, err := store.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Equal(t, originalID, got.ID)
	assert.Equal(t, originalCreatedAt, got.CreatedAt)
	assert.Len(t, got.Sessions, 2)
	assert.Equal(t, "Updated Day 1", got.Sessions[0].SessionName)
}

func TestPlanStore_AddSessions(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()
	userID := "user-1"

	plan := &domain.Plan{
		Sessions: []domain.PlanSession{
			{SessionName: "Day 1", Order: 0},
		},
	}
	require.NoError(t, store.Create(ctx, userID, plan))

	newSessions := []domain.PlanSession{
		{SessionName: "Day 2", Order: 1, Entries: []domain.PlanSessionEntry{
			{ExerciseName: "Bench", Order: 0},
		}},
	}
	err := store.AddSessions(ctx, userID, newSessions)
	require.NoError(t, err)

	got, err := store.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, got.Sessions, 2)
	assert.Equal(t, "Day 2", got.Sessions[1].SessionName)
	assert.NotEqual(t, uuid.Nil, got.Sessions[1].ID)
}

func TestPlanStore_DeleteSession(t *testing.T) {
	store := NewPlanRepository()
	ctx := context.Background()
	userID := "user-1"

	plan := &domain.Plan{
		Sessions: []domain.PlanSession{
			{SessionName: "Day 1", Order: 0},
			{SessionName: "Day 2", Order: 1},
		},
	}
	require.NoError(t, store.Create(ctx, userID, plan))

	got, err := store.GetByUserID(ctx, userID)
	require.NoError(t, err)
	sessionID := got.Sessions[0].ID.String()

	err = store.DeleteSession(ctx, userID, sessionID)
	require.NoError(t, err)

	got, err = store.GetByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, got.Sessions, 1)
	assert.Equal(t, "Day 2", got.Sessions[0].SessionName)

	// Delete nonexistent session
	err = store.DeleteSession(ctx, userID, uuid.New().String())
	assert.ErrorIs(t, err, domain.ErrNotFound)
}
