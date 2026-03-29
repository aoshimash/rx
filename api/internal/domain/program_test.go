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

func TestValidateProgram(t *testing.T) {
	validProgram := func() *Program {
		return &Program{
			ID:     uuid.New(),
			Name:   "Test Program",
			Status: ProgramStatusCreated,
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
							Fields:       map[string]interface{}{"sets": float64(3), "reps": float64(10)},
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
			Status: ProgramStatusCreated,
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

	t.Run("valid with ongoing status", func(t *testing.T) {
		p := validProgram()
		p.Status = ProgramStatusOngoing
		err := ValidateProgram(p)
		assert.NoError(t, err)
	})

	t.Run("valid with completed status", func(t *testing.T) {
		p := validProgram()
		p.Status = ProgramStatusCompleted
		err := ValidateProgram(p)
		assert.NoError(t, err)
	})

	t.Run("valid with cancelled status", func(t *testing.T) {
		p := validProgram()
		p.Status = ProgramStatusCancelled
		err := ValidateProgram(p)
		assert.NoError(t, err)
	})

	t.Run("invalid status", func(t *testing.T) {
		p := validProgram()
		p.Status = "invalid"
		err := ValidateProgram(p)
		assert.Error(t, err)
	})

	t.Run("old active status is invalid", func(t *testing.T) {
		p := validProgram()
		p.Status = "active"
		err := ValidateProgram(p)
		assert.Error(t, err)
	})
}

func TestValidateProgramStatusTransition(t *testing.T) {
	t.Run("created to ongoing is valid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCreated, ProgramStatusOngoing)
		assert.NoError(t, err)
	})

	t.Run("ongoing to completed is valid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusOngoing, ProgramStatusCompleted)
		assert.NoError(t, err)
	})

	t.Run("ongoing to cancelled is valid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusOngoing, ProgramStatusCancelled)
		assert.NoError(t, err)
	})

	t.Run("created to completed is invalid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCreated, ProgramStatusCompleted)
		assert.Error(t, err)
	})

	t.Run("completed to ongoing is invalid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCompleted, ProgramStatusOngoing)
		assert.Error(t, err)
	})

	t.Run("cancelled to ongoing is valid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCancelled, ProgramStatusOngoing)
		assert.NoError(t, err)
	})

	t.Run("cancelled to created is invalid", func(t *testing.T) {
		err := ValidateProgramStatusTransition(ProgramStatusCancelled, ProgramStatusCreated)
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
