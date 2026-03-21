package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func mustFormatDate(d DateOnly) string {
	return time.Time(d).Format("2006-01-02")
}

func TestRoundToIncrement(t *testing.T) {
	t.Run("rounds to 2.5kg increment", func(t *testing.T) {
		assert.Equal(t, 80.0, RoundToIncrement(81.0, 2.5))
		assert.Equal(t, 82.5, RoundToIncrement(82.0, 2.5))
		assert.Equal(t, 80.0, RoundToIncrement(79.0, 2.5))
	})

	t.Run("rounds to 5kg increment", func(t *testing.T) {
		assert.Equal(t, 80.0, RoundToIncrement(82.0, 5.0))
		assert.Equal(t, 85.0, RoundToIncrement(83.0, 5.0))
	})

	t.Run("rounds to 1kg increment", func(t *testing.T) {
		assert.Equal(t, 60.0, RoundToIncrement(60.4, 1.0))
		assert.Equal(t, 61.0, RoundToIncrement(60.5, 1.0))
	})

	t.Run("zero increment falls back to 0.1kg", func(t *testing.T) {
		assert.Equal(t, 60.2, RoundToIncrement(60.15, 0))
	})

	t.Run("negative increment falls back to 0.1kg", func(t *testing.T) {
		assert.Equal(t, 60.2, RoundToIncrement(60.15, -1))
	})

	t.Run("exact value unchanged", func(t *testing.T) {
		assert.Equal(t, 100.0, RoundToIncrement(100.0, 2.5))
	})
}

func TestConvertProgramToPlans(t *testing.T) {
	sets3 := 3
	reps5 := 5
	rpe8 := 8
	pct80 := 0.80
	pct70 := 0.70

	// Program with no session metadata (single session)
	program := &Program{
		ID:   uuid.New(),
		Name: "Strength Program",
		Entries: []ProgramEntry{
			{
				ID:           uuid.New(),
				Order:        0,
				ExerciseName: "Squat",
				Sets:         &sets3,
				Reps:         &reps5,
				RPE:          &rpe8,
				Percent1RM:   &pct80,
			},
			{
				ID:           uuid.New(),
				Order:        1,
				ExerciseName: "Bench Press",
				Sets:         &sets3,
				Reps:         &reps5,
				RPE:          &rpe8,
				Percent1RM:   &pct70,
			},
			{
				ID:           uuid.New(),
				Order:        2,
				ExerciseName: "Chin Up",
				Sets:         &sets3,
				Reps:         &reps5,
			},
		},
	}

	t.Run("single session - calculates weights from percent_1rm", func(t *testing.T) {
		input := &ConvertProgramToPlansInput{
			TargetWeights: map[string]float64{
				"Squat":       200.0,
				"Bench Press": 100.0,
			},
			LoadIncrements: map[string]float64{
				"Squat":       2.5,
				"Bench Press": 2.5,
			},
		}

		plans := ConvertProgramToPlans(program, input)

		require.Len(t, plans, 1)
		plan := plans[0]
		assert.Equal(t, "Strength Program", plan.Name)
		assert.Equal(t, &program.ID, plan.ProgramID)
		assert.Nil(t, plan.SessionName) // no session metadata
		require.Len(t, plan.Entries, 3)

		// Squat: 0.80 * 200 = 160.0, rounded to 2.5 → 160.0
		require.NotNil(t, plan.Entries[0].LoadKg)
		assert.Equal(t, 160.0, *plan.Entries[0].LoadKg)
		assert.Equal(t, "Squat", plan.Entries[0].ExerciseName)

		// Bench Press: 0.70 * 100 = 70.0, rounded to 2.5 → 70.0
		require.NotNil(t, plan.Entries[1].LoadKg)
		assert.Equal(t, 70.0, *plan.Entries[1].LoadKg)

		// Chin Up: no target weight provided → nil
		assert.Nil(t, plan.Entries[2].LoadKg)
	})

	t.Run("copies direct weight when no percent_1rm", func(t *testing.T) {
		input := &ConvertProgramToPlansInput{
			TargetWeights: map[string]float64{
				"Chin Up": 10.0,
			},
			LoadIncrements: map[string]float64{
				"Chin Up": 2.5,
			},
		}

		plans := ConvertProgramToPlans(program, input)

		require.Len(t, plans, 1)
		require.NotNil(t, plans[0].Entries[2].LoadKg)
		assert.Equal(t, 10.0, *plans[0].Entries[2].LoadKg)
	})

	t.Run("rounds to increment", func(t *testing.T) {
		input := &ConvertProgramToPlansInput{
			TargetWeights: map[string]float64{
				"Squat": 195.0,
			},
			LoadIncrements: map[string]float64{
				"Squat": 2.5,
			},
		}

		plans := ConvertProgramToPlans(program, input)

		require.Len(t, plans, 1)
		require.NotNil(t, plans[0].Entries[0].LoadKg)
		assert.Equal(t, 155.0, *plans[0].Entries[0].LoadKg)
	})

	t.Run("uses 0.1kg precision without increment", func(t *testing.T) {
		input := &ConvertProgramToPlansInput{
			TargetWeights: map[string]float64{
				"Squat": 195.0,
			},
		}

		plans := ConvertProgramToPlans(program, input)

		require.Len(t, plans, 1)
		require.NotNil(t, plans[0].Entries[0].LoadKg)
		assert.Equal(t, 156.0, *plans[0].Entries[0].LoadKg)
	})

	t.Run("custom plan name", func(t *testing.T) {
		input := &ConvertProgramToPlansInput{
			Name:          "Week 1 Plan",
			TargetWeights: map[string]float64{},
		}

		plans := ConvertProgramToPlans(program, input)

		require.Len(t, plans, 1)
		assert.Equal(t, "Week 1 Plan", plans[0].Name)
	})

	t.Run("copies sets reps rpe notes", func(t *testing.T) {
		input := &ConvertProgramToPlansInput{
			TargetWeights: map[string]float64{},
		}

		plans := ConvertProgramToPlans(program, input)

		require.Len(t, plans, 1)
		assert.Equal(t, &sets3, plans[0].Entries[0].Sets)
		assert.Equal(t, &reps5, plans[0].Entries[0].Reps)
		assert.Equal(t, &rpe8, plans[0].Entries[0].RPE)
	})

	t.Run("stores conversion metadata", func(t *testing.T) {
		input := &ConvertProgramToPlansInput{
			TargetWeights: map[string]float64{
				"Squat": 200.0,
			},
			LoadIncrements: map[string]float64{
				"Squat": 2.5,
			},
		}

		plans := ConvertProgramToPlans(program, input)

		require.Len(t, plans, 1)
		var meta map[string]interface{}
		err := json.Unmarshal(plans[0].Metadata, &meta)
		require.NoError(t, err)

		conversion, ok := meta["conversion"].(map[string]interface{})
		require.True(t, ok)
		tw, ok := conversion["target_weights"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, 200.0, tw["Squat"])
	})

	t.Run("multi-session groups entries by metadata.session", func(t *testing.T) {
		multiProgram := &Program{
			ID:   uuid.New(),
			Name: "Upper/Lower Split",
			Entries: []ProgramEntry{
				{
					Order:        0,
					ExerciseName: "Squat",
					Sets:         &sets3,
					Reps:         &reps5,
					Percent1RM:   &pct80,
					Metadata:     json.RawMessage(`{"session": "Lower"}`),
				},
				{
					Order:        1,
					ExerciseName: "Deadlift",
					Sets:         &sets3,
					Reps:         &reps5,
					Metadata:     json.RawMessage(`{"session": "Lower"}`),
				},
				{
					Order:        2,
					ExerciseName: "Bench Press",
					Sets:         &sets3,
					Reps:         &reps5,
					Percent1RM:   &pct70,
					Metadata:     json.RawMessage(`{"session": "Upper"}`),
				},
				{
					Order:        3,
					ExerciseName: "Row",
					Sets:         &sets3,
					Reps:         &reps5,
					Metadata:     json.RawMessage(`{"session": "Upper"}`),
				},
			},
		}

		input := &ConvertProgramToPlansInput{
			TargetWeights: map[string]float64{
				"Squat":       200.0,
				"Bench Press": 100.0,
			},
		}

		plans := ConvertProgramToPlans(multiProgram, input)

		require.Len(t, plans, 2)

		// First plan: Lower session
		assert.Equal(t, "Lower", plans[0].Name)
		require.NotNil(t, plans[0].SessionName)
		assert.Equal(t, "Lower", *plans[0].SessionName)
		require.Len(t, plans[0].Entries, 2)
		assert.Equal(t, "Squat", plans[0].Entries[0].ExerciseName)
		assert.Equal(t, "Deadlift", plans[0].Entries[1].ExerciseName)

		// Second plan: Upper session
		assert.Equal(t, "Upper", plans[1].Name)
		require.NotNil(t, plans[1].SessionName)
		assert.Equal(t, "Upper", *plans[1].SessionName)
		require.Len(t, plans[1].Entries, 2)
		assert.Equal(t, "Bench Press", plans[1].Entries[0].ExerciseName)
		assert.Equal(t, "Row", plans[1].Entries[1].ExerciseName)
	})

	t.Run("copies entry metadata", func(t *testing.T) {
		programWithMeta := &Program{
			ID:   uuid.New(),
			Name: "Meta Program",
			Entries: []ProgramEntry{
				{
					Order:        0,
					ExerciseName: "Squat",
					Metadata:     json.RawMessage(`{"week": 1}`),
				},
			},
		}

		input := &ConvertProgramToPlansInput{
			TargetWeights: map[string]float64{},
		}

		plans := ConvertProgramToPlans(programWithMeta, input)

		require.Len(t, plans, 1)
		assert.JSONEq(t, `{"week": 1}`, string(plans[0].Entries[0].Metadata))
	})

	t.Run("empty program produces empty plans", func(t *testing.T) {
		emptyProgram := &Program{
			ID:      uuid.New(),
			Name:    "Empty",
			Entries: []ProgramEntry{},
		}

		input := &ConvertProgramToPlansInput{
			TargetWeights: map[string]float64{},
		}

		plans := ConvertProgramToPlans(emptyProgram, input)
		assert.Empty(t, plans)
	})
}

func TestValidateProgramEntry(t *testing.T) {
	validEntry := func() *ProgramEntry {
		sets := 3
		reps := 10
		rpe := 8
		pct := 0.75
		return &ProgramEntry{
			ID:           uuid.New(),
			ProgramID:    uuid.New(),
			Order:        0,
			ExerciseName: "Squat",
			Sets:         &sets,
			Reps:         &reps,
			RPE:          &rpe,
			Percent1RM:   &pct,
		}
	}

	t.Run("valid entry", func(t *testing.T) {
		e := validEntry()
		err := ValidateProgramEntry(e)
		assert.NoError(t, err)
	})

	t.Run("nil entry", func(t *testing.T) {
		err := ValidateProgramEntry(nil)
		assert.Error(t, err)
	})

	t.Run("empty exercise_name", func(t *testing.T) {
		e := validEntry()
		e.ExerciseName = ""
		err := ValidateProgramEntry(e)
		assert.Error(t, err)
	})

	t.Run("negative order", func(t *testing.T) {
		e := validEntry()
		e.Order = -1
		err := ValidateProgramEntry(e)
		assert.Error(t, err)
	})

	t.Run("zero sets", func(t *testing.T) {
		e := validEntry()
		sets := 0
		e.Sets = &sets
		err := ValidateProgramEntry(e)
		assert.Error(t, err)
	})

	t.Run("zero reps", func(t *testing.T) {
		e := validEntry()
		reps := 0
		e.Reps = &reps
		err := ValidateProgramEntry(e)
		assert.Error(t, err)
	})

	t.Run("invalid RPE", func(t *testing.T) {
		e := validEntry()
		rpe := 11
		e.RPE = &rpe
		err := ValidateProgramEntry(e)
		assert.Error(t, err)
	})

	t.Run("percent_1rm out of range high", func(t *testing.T) {
		e := validEntry()
		pct := 1.5
		e.Percent1RM = &pct
		err := ValidateProgramEntry(e)
		assert.Error(t, err)
	})

	t.Run("percent_1rm out of range low", func(t *testing.T) {
		e := validEntry()
		pct := -0.1
		e.Percent1RM = &pct
		err := ValidateProgramEntry(e)
		assert.Error(t, err)
	})

	t.Run("percent_1rm at boundaries", func(t *testing.T) {
		e := validEntry()
		pct := 0.0
		e.Percent1RM = &pct
		assert.NoError(t, ValidateProgramEntry(e))

		pct = 1.0
		e.Percent1RM = &pct
		assert.NoError(t, ValidateProgramEntry(e))
	})

	t.Run("nil optional fields are valid", func(t *testing.T) {
		e := &ProgramEntry{
			ID:           uuid.New(),
			ProgramID:    uuid.New(),
			Order:        0,
			ExerciseName: "Squat",
		}
		err := ValidateProgramEntry(e)
		assert.NoError(t, err)
	})
}

func TestValidateProgram(t *testing.T) {
	validProgram := func() *Program {
		sets := 3
		reps := 10
		return &Program{
			ID:   uuid.New(),
			Name: "Test Program",
			Entries: []ProgramEntry{
				{
					ID:           uuid.New(),
					ProgramID:    uuid.New(),
					Order:        0,
					ExerciseName: "Squat",
					Sets:         &sets,
					Reps:         &reps,
				},
			},
		}
	}

	t.Run("valid program", func(t *testing.T) {
		p := validProgram()
		err := ValidateProgram(p)
		assert.NoError(t, err)
	})

	t.Run("nil program", func(t *testing.T) {
		err := ValidateProgram(nil)
		assert.Error(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		p := validProgram()
		p.Name = ""
		err := ValidateProgram(p)
		assert.Error(t, err)
	})

	t.Run("too many entries", func(t *testing.T) {
		p := validProgram()
		p.Entries = make([]ProgramEntry, 1001)
		for i := range p.Entries {
			p.Entries[i] = ProgramEntry{
				ID:           uuid.New(),
				ProgramID:    p.ID,
				Order:        i,
				ExerciseName: "Exercise",
			}
		}
		err := ValidateProgram(p)
		assert.Error(t, err)
	})

	t.Run("program with no entries is valid", func(t *testing.T) {
		p := &Program{
			ID:   uuid.New(),
			Name: "Empty Program",
		}
		err := ValidateProgram(p)
		assert.NoError(t, err)
	})
}

func TestDateOnly(t *testing.T) {
	t.Run("marshal to JSON", func(t *testing.T) {
		d := DateOnly(mustParseDate("2026-03-16"))
		data, err := json.Marshal(d)
		require.NoError(t, err)
		assert.Equal(t, `"2026-03-16"`, string(data))
	})

	t.Run("unmarshal from JSON", func(t *testing.T) {
		var d DateOnly
		err := json.Unmarshal([]byte(`"2026-03-16"`), &d)
		require.NoError(t, err)
		assert.Equal(t, "2026-03-16", mustFormatDate(d))
	})

	t.Run("unmarshal invalid format", func(t *testing.T) {
		var d DateOnly
		err := json.Unmarshal([]byte(`"invalid"`), &d)
		assert.Error(t, err)
	})
}
