package memory

import (
	"context"
	"testing"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgramMemoryStore_ExistsByName(t *testing.T) {
	store := NewProgramRepository()
	ctx := context.Background()

	t.Run("returns false when empty", func(t *testing.T) {
		exists, err := store.ExistsByName(ctx, "Test")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("returns true after create", func(t *testing.T) {
		p := &domain.Program{Name: "Test", Status: domain.ProgramStatusCreated}
		require.NoError(t, store.Create(ctx, p))

		exists, err := store.ExistsByName(ctx, "Test")
		require.NoError(t, err)
		assert.True(t, exists)
	})
}
