package memory

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestLog(t *testing.T, repo *logStore, programID *uuid.UUID, sessionName *string) *domain.Log {
	t.Helper()
	log := &domain.Log{
		ProgramID:   programID,
		SessionName: sessionName,
		PerformedAt: time.Now(),
		Entries: []domain.LogEntry{
			{ExerciseName: "Squat", Order: 0},
		},
	}
	err := repo.Create(context.Background(), log)
	require.NoError(t, err)
	return log
}

func TestLogStore_List_ProgramIDFilter(t *testing.T) {
	store := NewLogRepository().(*logStore)
	ctx := context.Background()

	programID1 := uuid.New()
	programID2 := uuid.New()

	// Create logs: 2 for program1, 1 for program2, 1 without program
	createTestLog(t, store, &programID1, nil)
	createTestLog(t, store, &programID1, nil)
	createTestLog(t, store, &programID2, nil)
	createTestLog(t, store, nil, nil)

	t.Run("no filter returns all", func(t *testing.T) {
		logs, _, _, err := store.List(ctx, nil, 100, "")
		require.NoError(t, err)
		assert.Len(t, logs, 4)
	})

	t.Run("filter by program1", func(t *testing.T) {
		logs, _, _, err := store.List(ctx, &programID1, 100, "")
		require.NoError(t, err)
		assert.Len(t, logs, 2)
		for _, l := range logs {
			assert.Equal(t, programID1, *l.ProgramID)
		}
	})

	t.Run("filter by program2", func(t *testing.T) {
		logs, _, _, err := store.List(ctx, &programID2, 100, "")
		require.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, programID2, *logs[0].ProgramID)
	})

	t.Run("filter by nonexistent program", func(t *testing.T) {
		nonexistent := uuid.New()
		logs, _, _, err := store.List(ctx, &nonexistent, 100, "")
		require.NoError(t, err)
		assert.Len(t, logs, 0)
	})
}

func TestLogStore_ListByPerformedAtRange_ProgramIDFilter(t *testing.T) {
	store := NewLogRepository().(*logStore)
	ctx := context.Background()

	programID := uuid.New()
	now := time.Now()

	// Create a log with program at a known time
	log1 := &domain.Log{
		ProgramID:   &programID,
		PerformedAt: now,
		Entries:     []domain.LogEntry{{ExerciseName: "Squat", Order: 0}},
	}
	require.NoError(t, store.Create(ctx, log1))

	// Create a log without program at the same time
	log2 := &domain.Log{
		PerformedAt: now,
		Entries:     []domain.LogEntry{{ExerciseName: "Bench", Order: 0}},
	}
	require.NoError(t, store.Create(ctx, log2))

	from := now.Add(-time.Hour)
	to := now.Add(time.Hour)

	t.Run("combined filter returns matching log", func(t *testing.T) {
		logs, _, _, err := store.ListByPerformedAtRange(ctx, &programID, &from, &to, 100, "")
		require.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, programID, *logs[0].ProgramID)
	})

	t.Run("time filter without program returns all in range", func(t *testing.T) {
		logs, _, _, err := store.ListByPerformedAtRange(ctx, nil, &from, &to, 100, "")
		require.NoError(t, err)
		assert.Len(t, logs, 2)
	})
}

func TestLogStore_CreateAndGetWithSetsAndPlanSnapshot(t *testing.T) {
	store := NewLogRepository().(*logStore)
	ctx := context.Background()

	snapshot := json.RawMessage(`{"sessions":[{"name":"Day 1"}]}`)

	log := &domain.Log{
		PerformedAt:  time.Now(),
		PlanSnapshot: snapshot,
		Entries: []domain.LogEntry{
			{
				ExerciseName: "Squat",
				Order:        0,
				Sets: []domain.LogSet{
					{
						SetNumber: 1,
						Fields:    map[string]interface{}{"weight_kg": float64(100), "reps": float64(5)},
					},
					{
						SetNumber: 2,
						Fields:    map[string]interface{}{"weight_kg": float64(105), "reps": float64(3)},
					},
				},
			},
		},
	}

	err := store.Create(ctx, log)
	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, log.ID)

	got, err := store.GetByID(ctx, log.ID)
	require.NoError(t, err)

	// Verify PlanSnapshot
	assert.JSONEq(t, string(snapshot), string(got.PlanSnapshot))

	// Verify Sets
	require.Len(t, got.Entries, 1)
	require.Len(t, got.Entries[0].Sets, 2)
	assert.Equal(t, 1, got.Entries[0].Sets[0].SetNumber)
	assert.Equal(t, float64(100), got.Entries[0].Sets[0].Fields["weight_kg"])
	assert.Equal(t, 2, got.Entries[0].Sets[1].SetNumber)
	assert.NotEqual(t, uuid.Nil, got.Entries[0].Sets[0].ID)
	assert.Equal(t, got.Entries[0].ID, got.Entries[0].Sets[0].EntryID)

	// Verify deep copy isolation
	got.PlanSnapshot = json.RawMessage(`{"modified":true}`)
	got.Entries[0].Sets[0].Fields["weight_kg"] = float64(999)

	got2, err := store.GetByID(ctx, log.ID)
	require.NoError(t, err)
	assert.JSONEq(t, string(snapshot), string(got2.PlanSnapshot))
	assert.Equal(t, float64(100), got2.Entries[0].Sets[0].Fields["weight_kg"])
}
