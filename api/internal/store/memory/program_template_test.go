package memory

import (
	"context"
	"testing"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgramTemplateStore_Update(t *testing.T) {
	store := NewProgramTemplateRepository()
	ctx := context.Background()

	tmpl := &domain.ProgramTemplate{
		Name: "Original",
		Entries: []domain.ProgramTemplateEntry{
			{Order: 0, ExerciseName: "Squat"},
		},
	}
	require.NoError(t, store.Create(ctx, tmpl))
	originalID := tmpl.ID
	originalCreatedAt := tmpl.CreatedAt

	tmpl.Name = "Updated"
	tmpl.Entries = []domain.ProgramTemplateEntry{
		{Order: 0, ExerciseName: "Bench Press"},
		{Order: 1, ExerciseName: "Deadlift"},
	}
	require.NoError(t, store.Update(ctx, tmpl))

	updated, err := store.GetByID(ctx, originalID)
	require.NoError(t, err)
	assert.Equal(t, "Updated", updated.Name)
	assert.Equal(t, originalID, updated.ID)
	assert.Equal(t, originalCreatedAt, updated.CreatedAt)
	assert.True(t, updated.UpdatedAt.After(originalCreatedAt))
	assert.Len(t, updated.Entries, 2)
	assert.Equal(t, "Bench Press", updated.Entries[0].ExerciseName)
}

func TestProgramTemplateStore_Update_NotFound(t *testing.T) {
	store := NewProgramTemplateRepository()
	ctx := context.Background()

	tmpl := &domain.ProgramTemplate{
		Name: "Nonexistent",
	}
	err := store.Update(ctx, tmpl)
	assert.Equal(t, domain.ErrNotFound, err)
}
