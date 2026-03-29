package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateProgramGroup(t *testing.T) {
	validGroup := func() *ProgramGroup {
		return &ProgramGroup{
			ID:        uuid.New(),
			ProgramID: uuid.New(),
			Name:      "Block 1",
			Order:     0,
		}
	}

	t.Run("valid group", func(t *testing.T) {
		g := validGroup()
		err := ValidateProgramGroup(g)
		assert.NoError(t, err)
	})

	t.Run("nil group", func(t *testing.T) {
		err := ValidateProgramGroup(nil)
		assert.Error(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		g := validGroup()
		g.Name = ""
		err := ValidateProgramGroup(g)
		assert.Error(t, err)
	})

	t.Run("name too long", func(t *testing.T) {
		g := validGroup()
		g.Name = strings.Repeat("a", 201)
		err := ValidateProgramGroup(g)
		assert.Error(t, err)
	})

	t.Run("negative order", func(t *testing.T) {
		g := validGroup()
		g.Order = -1
		err := ValidateProgramGroup(g)
		assert.Error(t, err)
	})

	t.Run("notes too long", func(t *testing.T) {
		g := validGroup()
		longNotes := strings.Repeat("a", 5001)
		g.Notes = &longNotes
		err := ValidateProgramGroup(g)
		assert.Error(t, err)
	})

	t.Run("valid with notes", func(t *testing.T) {
		g := validGroup()
		notes := "Some notes"
		g.Notes = &notes
		err := ValidateProgramGroup(g)
		assert.NoError(t, err)
	})
}

func TestValidateGroupDepths(t *testing.T) {
	t.Run("depth 0 valid", func(t *testing.T) {
		groups := []ProgramGroup{
			{ID: uuid.New(), Name: "Block 1", Order: 0},
		}
		err := ValidateGroupDepths(groups)
		assert.NoError(t, err)
	})

	t.Run("depth 1 valid", func(t *testing.T) {
		parentID := uuid.New()
		groups := []ProgramGroup{
			{ID: parentID, Name: "Block 1", Order: 0},
			{ID: uuid.New(), ParentGroupID: &parentID, Name: "Week 1", Order: 0},
		}
		err := ValidateGroupDepths(groups)
		assert.NoError(t, err)
	})

	t.Run("depth 2 rejected", func(t *testing.T) {
		id1 := uuid.New()
		id2 := uuid.New()
		groups := []ProgramGroup{
			{ID: id1, Name: "Level 0", Order: 0},
			{ID: id2, ParentGroupID: &id1, Name: "Level 1", Order: 0},
			{ID: uuid.New(), ParentGroupID: &id2, Name: "Level 2", Order: 0},
		}
		err := ValidateGroupDepths(groups)
		assert.Error(t, err)
	})

	t.Run("empty groups valid", func(t *testing.T) {
		err := ValidateGroupDepths([]ProgramGroup{})
		assert.NoError(t, err)
	})

	t.Run("missing parent rejected", func(t *testing.T) {
		missingID := uuid.New()
		groups := []ProgramGroup{
			{ID: uuid.New(), ParentGroupID: &missingID, Name: "Orphan", Order: 0},
		}
		err := ValidateGroupDepths(groups)
		assert.Error(t, err)
	})
}
