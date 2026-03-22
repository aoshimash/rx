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
		assert.Equal(t, 80.0, RoundToIncrement(80.0, 2.5))
	})
}

func TestGenerateProgramFromTemplate(t *testing.T) {
	sets3 := 3
	reps5 := 5
	rpe8 := 8
	pct80 := 0.80
	pct70 := 0.70

	// Template with no session metadata (single session)
	tmpl := &ProgramTemplate{
		ID:   uuid.New(),
		Name: "Strength Program",
		Entries: []ProgramTemplateEntry{
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
		input := &GenerateProgramInput{
			TargetWeights: map[string]float64{
				"Squat":       200.0,
				"Bench Press": 100.0,
			},
			LoadIncrements: map[string]float64{
				"Squat":       2.5,
				"Bench Press": 2.5,
			},
		}

		program := GenerateProgramFromTemplate(tmpl, input)

		require.Len(t, program.Sessions, 1)
		session := program.Sessions[0]
		assert.Equal(t, "Strength Program", session.SessionName)
		require.Len(t, session.Entries, 3)

		// Squat: 0.80 * 200 = 160.0, rounded to 2.5 → 160.0
		require.NotNil(t, session.Entries[0].LoadKg)
		assert.Equal(t, 160.0, *session.Entries[0].LoadKg)
		assert.Equal(t, "Squat", session.Entries[0].ExerciseName)

		// Bench Press: 0.70 * 100 = 70.0, rounded to 2.5 → 70.0
		require.NotNil(t, session.Entries[1].LoadKg)
		assert.Equal(t, 70.0, *session.Entries[1].LoadKg)

		// Chin Up: no target weight provided → nil
		assert.Nil(t, session.Entries[2].LoadKg)
	})

	t.Run("copies direct weight when no percent_1rm", func(t *testing.T) {
		input := &GenerateProgramInput{
			TargetWeights: map[string]float64{
				"Chin Up": 10.0,
			},
			LoadIncrements: map[string]float64{
				"Chin Up": 2.5,
			},
		}

		program := GenerateProgramFromTemplate(tmpl, input)

		require.Len(t, program.Sessions, 1)
		require.NotNil(t, program.Sessions[0].Entries[2].LoadKg)
		assert.Equal(t, 10.0, *program.Sessions[0].Entries[2].LoadKg)
	})

	t.Run("rounds to increment", func(t *testing.T) {
		input := &GenerateProgramInput{
			TargetWeights: map[string]float64{
				"Squat": 195.0,
			},
			LoadIncrements: map[string]float64{
				"Squat": 2.5,
			},
		}

		program := GenerateProgramFromTemplate(tmpl, input)

		require.Len(t, program.Sessions, 1)
		require.NotNil(t, program.Sessions[0].Entries[0].LoadKg)
		assert.Equal(t, 155.0, *program.Sessions[0].Entries[0].LoadKg)
	})

	t.Run("uses 0.1kg precision without increment", func(t *testing.T) {
		input := &GenerateProgramInput{
			TargetWeights: map[string]float64{
				"Squat": 195.0,
			},
		}

		program := GenerateProgramFromTemplate(tmpl, input)

		require.Len(t, program.Sessions, 1)
		require.NotNil(t, program.Sessions[0].Entries[0].LoadKg)
		assert.Equal(t, 156.0, *program.Sessions[0].Entries[0].LoadKg)
	})

	t.Run("custom program name", func(t *testing.T) {
		input := &GenerateProgramInput{
			Name:          "Week 1",
			TargetWeights: map[string]float64{},
		}

		program := GenerateProgramFromTemplate(tmpl, input)

		assert.Equal(t, "Week 1", program.Name)
	})

	t.Run("copies sets reps rpe", func(t *testing.T) {
		input := &GenerateProgramInput{
			TargetWeights: map[string]float64{},
		}

		program := GenerateProgramFromTemplate(tmpl, input)

		require.Len(t, program.Sessions, 1)
		assert.Equal(t, &sets3, program.Sessions[0].Entries[0].Sets)
		assert.Equal(t, &reps5, program.Sessions[0].Entries[0].Reps)
		assert.Equal(t, &rpe8, program.Sessions[0].Entries[0].RPE)
	})

	t.Run("stores generation metadata", func(t *testing.T) {
		input := &GenerateProgramInput{
			TargetWeights: map[string]float64{
				"Squat": 200.0,
			},
			LoadIncrements: map[string]float64{
				"Squat": 2.5,
			},
		}

		program := GenerateProgramFromTemplate(tmpl, input)

		var meta map[string]interface{}
		err := json.Unmarshal(program.Metadata, &meta)
		require.NoError(t, err)

		generation, ok := meta["generation"].(map[string]interface{})
		require.True(t, ok)
		tw, ok := generation["target_weights"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, 200.0, tw["Squat"])
	})

	t.Run("multi-session groups entries by metadata.session", func(t *testing.T) {
		multiTmpl := &ProgramTemplate{
			ID:   uuid.New(),
			Name: "Upper/Lower Split",
			Entries: []ProgramTemplateEntry{
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

		input := &GenerateProgramInput{
			TargetWeights: map[string]float64{
				"Squat":       200.0,
				"Bench Press": 100.0,
			},
		}

		program := GenerateProgramFromTemplate(multiTmpl, input)

		require.Len(t, program.Sessions, 2)

		// First session: Lower
		assert.Equal(t, "Lower", program.Sessions[0].SessionName)
		require.Len(t, program.Sessions[0].Entries, 2)
		assert.Equal(t, "Squat", program.Sessions[0].Entries[0].ExerciseName)
		assert.Equal(t, "Deadlift", program.Sessions[0].Entries[1].ExerciseName)

		// Second session: Upper
		assert.Equal(t, "Upper", program.Sessions[1].SessionName)
		require.Len(t, program.Sessions[1].Entries, 2)
		assert.Equal(t, "Bench Press", program.Sessions[1].Entries[0].ExerciseName)
		assert.Equal(t, "Row", program.Sessions[1].Entries[1].ExerciseName)
	})

	t.Run("copies entry metadata", func(t *testing.T) {
		tmplWithMeta := &ProgramTemplate{
			ID:   uuid.New(),
			Name: "Meta Program",
			Entries: []ProgramTemplateEntry{
				{
					Order:        0,
					ExerciseName: "Squat",
					Metadata:     json.RawMessage(`{"week": 1}`),
				},
			},
		}

		input := &GenerateProgramInput{
			TargetWeights: map[string]float64{},
		}

		program := GenerateProgramFromTemplate(tmplWithMeta, input)

		require.Len(t, program.Sessions, 1)
		assert.JSONEq(t, `{"week": 1}`, string(program.Sessions[0].Entries[0].Metadata))
	})

	t.Run("empty template produces empty sessions", func(t *testing.T) {
		emptyTmpl := &ProgramTemplate{
			ID:      uuid.New(),
			Name:    "Empty",
			Entries: []ProgramTemplateEntry{},
		}

		input := &GenerateProgramInput{
			TargetWeights: map[string]float64{},
		}

		program := GenerateProgramFromTemplate(emptyTmpl, input)
		assert.Empty(t, program.Sessions)
	})
}

func TestValidateProgramTemplateEntry(t *testing.T) {
	validEntry := func() *ProgramTemplateEntry {
		sets := 3
		reps := 10
		rpe := 8
		pct := 0.75
		return &ProgramTemplateEntry{
			ID:                uuid.New(),
			ProgramTemplateID: uuid.New(),
			Order:             0,
			ExerciseName:      "Squat",
			Sets:              &sets,
			Reps:              &reps,
			RPE:               &rpe,
			Percent1RM:        &pct,
		}
	}

	t.Run("valid entry", func(t *testing.T) {
		e := validEntry()
		err := ValidateProgramTemplateEntry(e)
		assert.NoError(t, err)
	})

	t.Run("nil entry", func(t *testing.T) {
		err := ValidateProgramTemplateEntry(nil)
		assert.Error(t, err)
	})

	t.Run("empty exercise_name", func(t *testing.T) {
		e := validEntry()
		e.ExerciseName = ""
		err := ValidateProgramTemplateEntry(e)
		assert.Error(t, err)
	})

	t.Run("negative order", func(t *testing.T) {
		e := validEntry()
		e.Order = -1
		err := ValidateProgramTemplateEntry(e)
		assert.Error(t, err)
	})

	t.Run("zero sets", func(t *testing.T) {
		e := validEntry()
		sets := 0
		e.Sets = &sets
		err := ValidateProgramTemplateEntry(e)
		assert.Error(t, err)
	})

	t.Run("zero reps", func(t *testing.T) {
		e := validEntry()
		reps := 0
		e.Reps = &reps
		err := ValidateProgramTemplateEntry(e)
		assert.Error(t, err)
	})

	t.Run("invalid RPE", func(t *testing.T) {
		e := validEntry()
		rpe := 11
		e.RPE = &rpe
		err := ValidateProgramTemplateEntry(e)
		assert.Error(t, err)
	})

	t.Run("percent_1rm out of range high", func(t *testing.T) {
		e := validEntry()
		pct := 1.5
		e.Percent1RM = &pct
		err := ValidateProgramTemplateEntry(e)
		assert.Error(t, err)
	})

	t.Run("percent_1rm out of range low", func(t *testing.T) {
		e := validEntry()
		pct := -0.1
		e.Percent1RM = &pct
		err := ValidateProgramTemplateEntry(e)
		assert.Error(t, err)
	})

	t.Run("percent_1rm at boundaries", func(t *testing.T) {
		e := validEntry()
		pct := 0.0
		e.Percent1RM = &pct
		assert.NoError(t, ValidateProgramTemplateEntry(e))

		pct = 1.0
		e.Percent1RM = &pct
		assert.NoError(t, ValidateProgramTemplateEntry(e))
	})

	t.Run("nil optional fields are valid", func(t *testing.T) {
		e := &ProgramTemplateEntry{
			ID:           uuid.New(),
			Order:        0,
			ExerciseName: "Squat",
		}
		err := ValidateProgramTemplateEntry(e)
		assert.NoError(t, err)
	})
}

func TestValidateProgram(t *testing.T) {
	validProgram := func() *Program {
		sets := 3
		reps := 10
		return &Program{
			ID:     uuid.New(),
			Name:   "Test Program",
			Status: ProgramStatusActive,
			Sessions: []ProgramSession{
				{
					ID:          uuid.New(),
					ProgramID:   uuid.New(),
					SessionName: "Day 1",
					Order:       0,
					Entries: []ProgramSessionEntry{
						{
							ID:           uuid.New(),
							Order:        0,
							ExerciseName: "Squat",
							Sets:         &sets,
							Reps:         &reps,
						},
					},
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

	t.Run("program with no sessions is invalid", func(t *testing.T) {
		p := &Program{
			ID:     uuid.New(),
			Name:   "Empty Program",
			Status: ProgramStatusActive,
		}
		err := ValidateProgram(p)
		assert.Error(t, err)
	})

	t.Run("duplicate session_name is invalid", func(t *testing.T) {
		p := validProgram()
		p.Sessions = append(p.Sessions, ProgramSession{
			ID:          uuid.New(),
			ProgramID:   p.ID,
			SessionName: "Day 1", // same session_name → duplicate
			Order:       1,
			Entries: []ProgramSessionEntry{
				{
					ID:           uuid.New(),
					Order:        0,
					ExerciseName: "Bench Press",
				},
			},
		})
		err := ValidateProgram(p)
		assert.Error(t, err)
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
