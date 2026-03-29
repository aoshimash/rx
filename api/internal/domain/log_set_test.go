package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validLogSet() *LogSet {
	return &LogSet{
		ID:        uuid.New(),
		EntryID:   uuid.New(),
		SetNumber: 1,
		Fields:    map[string]interface{}{"reps": 5.0, "load_kg": 100.0},
	}
}

func TestValidateLogSet(t *testing.T) {
	t.Run("valid set", func(t *testing.T) {
		s := validLogSet()
		err := ValidateLogSet(s)
		require.NoError(t, err)
	})

	t.Run("nil set", func(t *testing.T) {
		err := ValidateLogSet(nil)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "log_set", ve.Field)
	})

	t.Run("zero set_number", func(t *testing.T) {
		s := validLogSet()
		s.SetNumber = 0
		err := ValidateLogSet(s)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "set_number", ve.Field)
	})

	t.Run("nil fields", func(t *testing.T) {
		s := validLogSet()
		s.Fields = nil
		err := ValidateLogSet(s)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "fields", ve.Field)
	})

	t.Run("notes too long", func(t *testing.T) {
		s := validLogSet()
		longNotes := strings.Repeat("a", 2001)
		s.Notes = &longNotes
		err := ValidateLogSet(s)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "notes", ve.Field)
	})

	t.Run("video_url too long", func(t *testing.T) {
		s := validLogSet()
		longURL := strings.Repeat("a", 2001)
		s.VideoURL = &longURL
		err := ValidateLogSet(s)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "video_url", ve.Field)
	})

	t.Run("empty video_url", func(t *testing.T) {
		s := validLogSet()
		empty := ""
		s.VideoURL = &empty
		err := ValidateLogSet(s)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "video_url", ve.Field)
	})
}

func TestValidateLogEntry_WithSets(t *testing.T) {
	t.Run("entry with valid sets", func(t *testing.T) {
		e := &LogEntry{
			ID:           uuid.New(),
			LogID:        uuid.New(),
			Order:        0,
			ExerciseName: "Bench Press",
			Fields:       map[string]interface{}{},
			Sets: []LogSet{
				{
					ID:        uuid.New(),
					EntryID:   uuid.New(),
					SetNumber: 1,
					Fields:    map[string]interface{}{"reps": 5.0},
				},
			},
		}
		err := ValidateLogEntry(e)
		require.NoError(t, err)
	})

	t.Run("entry with invalid set", func(t *testing.T) {
		e := &LogEntry{
			ID:           uuid.New(),
			LogID:        uuid.New(),
			Order:        0,
			ExerciseName: "Bench Press",
			Fields:       map[string]interface{}{},
			Sets: []LogSet{
				{
					ID:        uuid.New(),
					EntryID:   uuid.New(),
					SetNumber: 0, // invalid
					Fields:    map[string]interface{}{},
				},
			},
		}
		err := ValidateLogEntry(e)
		require.Error(t, err)
		var ve *ValidationError
		require.ErrorAs(t, err, &ve)
		assert.Equal(t, "sets[0]", ve.Field)
	})
}
