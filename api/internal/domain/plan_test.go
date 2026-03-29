package domain

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validPlan() *Plan {
	planID := uuid.New()
	sessionID := uuid.New()
	return &Plan{
		ID: planID,
		Sessions: []PlanSession{
			{
				ID:          sessionID,
				PlanID:      planID,
				SessionName: "Day 1",
				Order:       0,
				Entries: []PlanSessionEntry{
					{
						ID:           uuid.New(),
						SessionID:    sessionID,
						Order:        0,
						ExerciseName: "Squat",
						Fields:       map[string]interface{}{"sets": 3.0, "reps": 5.0},
					},
				},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func TestValidatePlan(t *testing.T) {
	t.Run("valid plan", func(t *testing.T) {
		p := validPlan()
		err := ValidatePlan(p)
		require.NoError(t, err)
	})

	t.Run("nil plan", func(t *testing.T) {
		err := ValidatePlan(nil)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "plan", ve.Field)
	})

	t.Run("empty sessions is valid", func(t *testing.T) {
		p := validPlan()
		p.Sessions = nil
		err := ValidatePlan(p)
		require.NoError(t, err)
	})

	t.Run("name too long", func(t *testing.T) {
		p := validPlan()
		longName := strings.Repeat("a", 201)
		p.Name = &longName
		err := ValidatePlan(p)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "name", ve.Field)
	})

	t.Run("session with empty name", func(t *testing.T) {
		p := validPlan()
		p.Sessions[0].SessionName = ""
		err := ValidatePlan(p)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "sessions[0]", ve.Field)
	})

	t.Run("entry with empty exercise name", func(t *testing.T) {
		p := validPlan()
		p.Sessions[0].Entries[0].ExerciseName = ""
		err := ValidatePlan(p)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "sessions[0]", ve.Field)
	})

	t.Run("too many sessions", func(t *testing.T) {
		p := validPlan()
		p.Sessions = make([]PlanSession, 201)
		for i := range p.Sessions {
			p.Sessions[i] = PlanSession{
				ID:          uuid.New(),
				PlanID:      p.ID,
				SessionName: "Session",
				Order:       i,
			}
		}
		err := ValidatePlan(p)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "sessions", ve.Field)
	})
}
