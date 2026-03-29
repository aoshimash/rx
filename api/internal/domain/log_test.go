package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateLogEntry(t *testing.T) {
	validEntry := func() *LogEntry {
		return &LogEntry{
			ID:           uuid.New(),
			LogID:        uuid.New(),
			Order:        0,
			ExerciseName: "Bench Press",
			Fields: map[string]interface{}{
				"sets":    float64(3),
				"reps":    float64(10),
				"load_kg": float64(60),
			},
		}
	}

	t.Run("valid entry", func(t *testing.T) {
		e := validEntry()
		err := ValidateLogEntry(e)
		assert.NoError(t, err)
	})

	t.Run("nil entry", func(t *testing.T) {
		err := ValidateLogEntry(nil)
		assert.Error(t, err)
	})

	t.Run("empty exercise_name", func(t *testing.T) {
		e := validEntry()
		e.ExerciseName = ""
		err := ValidateLogEntry(e)
		assert.Error(t, err)
	})

	t.Run("negative order", func(t *testing.T) {
		e := validEntry()
		e.Order = -1
		err := ValidateLogEntry(e)
		assert.Error(t, err)
	})

	t.Run("zero sets", func(t *testing.T) {
		e := validEntry()
		e.Fields["sets"] = float64(0)
		err := ValidateLogEntry(e)
		assert.Error(t, err)
	})

	t.Run("zero reps", func(t *testing.T) {
		e := validEntry()
		e.Fields["reps"] = float64(0)
		err := ValidateLogEntry(e)
		assert.Error(t, err)
	})

	t.Run("negative load_kg", func(t *testing.T) {
		e := validEntry()
		e.Fields["load_kg"] = float64(-1)
		err := ValidateLogEntry(e)
		assert.Error(t, err)
	})

	t.Run("load_kg rounded to 0.1", func(t *testing.T) {
		e := validEntry()
		e.Fields["load_kg"] = float64(60.15)
		err := ValidateLogEntry(e)
		require.NoError(t, err)
		assert.Equal(t, float64(60.2), e.Fields["load_kg"])
	})

	t.Run("nil optional fields are valid", func(t *testing.T) {
		e := &LogEntry{
			ID:           uuid.New(),
			LogID:        uuid.New(),
			Order:        0,
			ExerciseName: "Squat",
		}
		err := ValidateLogEntry(e)
		assert.NoError(t, err)
	})

	t.Run("valid video_object_key", func(t *testing.T) {
		e := validEntry()
		key := "videos/test.mp4"
		e.VideoObjectKey = &key
		err := ValidateLogEntry(e)
		assert.NoError(t, err)
	})

	t.Run("empty video_object_key", func(t *testing.T) {
		e := validEntry()
		key := ""
		e.VideoObjectKey = &key
		err := ValidateLogEntry(e)
		assert.Error(t, err)
	})

	t.Run("fields with extra data is accepted", func(t *testing.T) {
		e := validEntry()
		e.Fields["tempo"] = "3-1-1-0"
		err := ValidateLogEntry(e)
		assert.NoError(t, err)
	})

	t.Run("valid started_at and finished_at", func(t *testing.T) {
		e := validEntry()
		start := time.Now().Add(-2 * time.Hour)
		finish := time.Now().Add(-1 * time.Hour)
		e.StartedAt = &start
		e.FinishedAt = &finish
		err := ValidateLogEntry(e)
		assert.NoError(t, err)
	})

	t.Run("only started_at is valid", func(t *testing.T) {
		e := validEntry()
		start := time.Now().Add(-2 * time.Hour)
		e.StartedAt = &start
		err := ValidateLogEntry(e)
		assert.NoError(t, err)
	})

	t.Run("only finished_at is valid", func(t *testing.T) {
		e := validEntry()
		finish := time.Now().Add(-1 * time.Hour)
		e.FinishedAt = &finish
		err := ValidateLogEntry(e)
		assert.NoError(t, err)
	})

	t.Run("started_at after finished_at", func(t *testing.T) {
		e := validEntry()
		start := time.Now().Add(-1 * time.Hour)
		finish := time.Now().Add(-2 * time.Hour)
		e.StartedAt = &start
		e.FinishedAt = &finish
		err := ValidateLogEntry(e)
		assert.Error(t, err)
	})
}

func TestValidateLog(t *testing.T) {
	validLog := func() *Log {
		return &Log{
			ID:          uuid.New(),
			PerformedAt: time.Now().Add(-1 * time.Hour),
			Entries: []LogEntry{
				{
					ID:           uuid.New(),
					LogID:        uuid.New(),
					Order:        0,
					ExerciseName: "Bench Press",
					Fields: map[string]interface{}{
						"sets":    float64(3),
						"reps":    float64(10),
						"load_kg": float64(60),
					},
				},
			},
		}
	}

	t.Run("valid log", func(t *testing.T) {
		l := validLog()
		err := ValidateLog(l)
		assert.NoError(t, err)
	})

	t.Run("nil log", func(t *testing.T) {
		err := ValidateLog(nil)
		assert.Error(t, err)
	})

	t.Run("future performed_at", func(t *testing.T) {
		l := validLog()
		l.PerformedAt = time.Now().Add(1 * time.Hour)
		err := ValidateLog(l)
		assert.Error(t, err)
	})

	t.Run("no entries", func(t *testing.T) {
		l := validLog()
		l.Entries = nil
		err := ValidateLog(l)
		assert.Error(t, err)
	})

	t.Run("too many entries", func(t *testing.T) {
		l := validLog()
		l.Entries = make([]LogEntry, 501)
		for i := range l.Entries {
			l.Entries[i] = LogEntry{
				ID:           uuid.New(),
				LogID:        l.ID,
				Order:        i,
				ExerciseName: "Exercise",
			}
		}
		err := ValidateLog(l)
		assert.Error(t, err)
	})

	t.Run("log with program_id and session_name", func(t *testing.T) {
		l := validLog()
		programID := uuid.New()
		sessionName := "Week 1 Day 1"
		l.ProgramID = &programID
		l.SessionName = &sessionName
		err := ValidateLog(l)
		assert.NoError(t, err)
	})

	t.Run("log with metadata", func(t *testing.T) {
		l := validLog()
		l.Metadata = json.RawMessage(`{"body_weight_kg": 75.5, "fatigue_level": 3}`)
		_ = l // metadata field still exists on Log, so this is valid
		err := ValidateLog(l)
		assert.NoError(t, err)
	})

	t.Run("valid started_at and finished_at", func(t *testing.T) {
		l := validLog()
		start := time.Now().Add(-2 * time.Hour)
		finish := time.Now().Add(-1 * time.Hour)
		l.StartedAt = &start
		l.FinishedAt = &finish
		err := ValidateLog(l)
		assert.NoError(t, err)
	})

	t.Run("only started_at is valid", func(t *testing.T) {
		l := validLog()
		start := time.Now().Add(-2 * time.Hour)
		l.StartedAt = &start
		err := ValidateLog(l)
		assert.NoError(t, err)
	})

	t.Run("only finished_at is valid", func(t *testing.T) {
		l := validLog()
		finish := time.Now().Add(-1 * time.Hour)
		l.FinishedAt = &finish
		err := ValidateLog(l)
		assert.NoError(t, err)
	})

	t.Run("started_at after finished_at", func(t *testing.T) {
		l := validLog()
		start := time.Now().Add(-1 * time.Hour)
		finish := time.Now().Add(-2 * time.Hour)
		l.StartedAt = &start
		l.FinishedAt = &finish
		err := ValidateLog(l)
		assert.Error(t, err)
	})
}
