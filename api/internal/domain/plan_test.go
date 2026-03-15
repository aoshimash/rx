package domain

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePlanEntry(t *testing.T) {
	validEntry := func() *PlanEntry {
		sets := 3
		reps := 10
		loadKg := 60.0
		rpe := 8
		return &PlanEntry{
			ID:           uuid.New(),
			PlanID:       uuid.New(),
			Order:        0,
			ExerciseName: "Bench Press",
			Sets:         &sets,
			Reps:         &reps,
			LoadKg:       &loadKg,
			RPE:          &rpe,
		}
	}

	t.Run("valid entry", func(t *testing.T) {
		e := validEntry()
		err := ValidatePlanEntry(e)
		assert.NoError(t, err)
	})

	t.Run("nil entry", func(t *testing.T) {
		err := ValidatePlanEntry(nil)
		assert.Error(t, err)
	})

	t.Run("empty exercise_name", func(t *testing.T) {
		e := validEntry()
		e.ExerciseName = ""
		err := ValidatePlanEntry(e)
		assert.Error(t, err)
	})

	t.Run("negative order", func(t *testing.T) {
		e := validEntry()
		e.Order = -1
		err := ValidatePlanEntry(e)
		assert.Error(t, err)
	})

	t.Run("zero sets", func(t *testing.T) {
		e := validEntry()
		sets := 0
		e.Sets = &sets
		err := ValidatePlanEntry(e)
		assert.Error(t, err)
	})

	t.Run("zero reps", func(t *testing.T) {
		e := validEntry()
		reps := 0
		e.Reps = &reps
		err := ValidatePlanEntry(e)
		assert.Error(t, err)
	})

	t.Run("negative load_kg", func(t *testing.T) {
		e := validEntry()
		loadKg := -1.0
		e.LoadKg = &loadKg
		err := ValidatePlanEntry(e)
		assert.Error(t, err)
	})

	t.Run("load_kg rounded to 0.1", func(t *testing.T) {
		e := validEntry()
		loadKg := 60.15
		e.LoadKg = &loadKg
		err := ValidatePlanEntry(e)
		require.NoError(t, err)
		assert.Equal(t, 60.2, *e.LoadKg)
	})

	t.Run("invalid RPE", func(t *testing.T) {
		e := validEntry()
		rpe := 11
		e.RPE = &rpe
		err := ValidatePlanEntry(e)
		assert.Error(t, err)
	})

	t.Run("nil optional fields are valid", func(t *testing.T) {
		e := &PlanEntry{
			ID:           uuid.New(),
			PlanID:       uuid.New(),
			Order:        0,
			ExerciseName: "Squat",
		}
		err := ValidatePlanEntry(e)
		assert.NoError(t, err)
	})

	t.Run("metadata is accepted", func(t *testing.T) {
		e := validEntry()
		e.Metadata = json.RawMessage(`{"week": 1, "day": "Monday"}`)
		err := ValidatePlanEntry(e)
		assert.NoError(t, err)
	})
}

func TestValidatePlan(t *testing.T) {
	validPlan := func() *Plan {
		sets := 3
		reps := 10
		return &Plan{
			ID:   uuid.New(),
			Name: "Test Plan",
			Entries: []PlanEntry{
				{
					ID:           uuid.New(),
					PlanID:       uuid.New(),
					Order:        0,
					ExerciseName: "Bench Press",
					Sets:         &sets,
					Reps:         &reps,
				},
			},
		}
	}

	t.Run("valid plan", func(t *testing.T) {
		p := validPlan()
		err := ValidatePlan(p)
		assert.NoError(t, err)
	})

	t.Run("nil plan", func(t *testing.T) {
		err := ValidatePlan(nil)
		assert.Error(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		p := validPlan()
		p.Name = ""
		err := ValidatePlan(p)
		assert.Error(t, err)
	})

	t.Run("too many entries", func(t *testing.T) {
		p := validPlan()
		p.Entries = make([]PlanEntry, 1001)
		for i := range p.Entries {
			p.Entries[i] = PlanEntry{
				ID:           uuid.New(),
				PlanID:       p.ID,
				Order:        i,
				ExerciseName: "Exercise",
			}
		}
		err := ValidatePlan(p)
		assert.Error(t, err)
	})

	t.Run("plan with no entries is valid", func(t *testing.T) {
		p := &Plan{
			ID:   uuid.New(),
			Name: "Empty Plan",
		}
		err := ValidatePlan(p)
		assert.NoError(t, err)
	})

	t.Run("plan with metadata", func(t *testing.T) {
		p := validPlan()
		p.Metadata = json.RawMessage(`{"source": "ai"}`)
		err := ValidatePlan(p)
		assert.NoError(t, err)
	})
}
