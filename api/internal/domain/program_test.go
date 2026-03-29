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
			ID:   uuid.New(),
			Name: "Test Program",
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
			ID:   uuid.New(),
			Name: "Empty Program",
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

	t.Run("valid program with groups", func(t *testing.T) {
		p := validProgram()
		blockID := uuid.New()
		weekID := uuid.New()
		p.Groups = []ProgramGroup{
			{
				ID:        blockID,
				ProgramID: p.ID,
				Name:      "Block 1",
				Order:     0,
			},
			{
				ID:            weekID,
				ProgramID:     p.ID,
				ParentGroupID: &blockID,
				Name:          "Week 1",
				Order:         0,
			},
		}
		p.Sessions[0].GroupID = &weekID
		err := ValidateProgram(p)
		assert.NoError(t, err)
	})

	t.Run("invalid group in program", func(t *testing.T) {
		p := validProgram()
		p.Groups = []ProgramGroup{
			{
				ID:        uuid.New(),
				ProgramID: p.ID,
				Name:      "", // invalid: empty name
				Order:     0,
			},
		}
		err := ValidateProgram(p)
		assert.Error(t, err)
	})

	t.Run("group depth exceeded in program", func(t *testing.T) {
		p := validProgram()
		id1 := uuid.New()
		id2 := uuid.New()
		id3 := uuid.New()
		p.Groups = []ProgramGroup{
			{ID: id1, ProgramID: p.ID, Name: "Level 0", Order: 0},
			{ID: id2, ProgramID: p.ID, ParentGroupID: &id1, Name: "Level 1", Order: 0},
			{ID: id3, ProgramID: p.ID, ParentGroupID: &id2, Name: "Level 2", Order: 0},
		}
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
