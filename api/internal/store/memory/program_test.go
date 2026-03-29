package memory

import (
	"context"
	"testing"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProgramMemoryStore_ExistsByName(t *testing.T) {
	tests := []struct {
		name            string
		programToCreate *domain.Program
		queryName       string
		wantExists      bool
	}{
		{
			name:            "returns false when empty",
			programToCreate: nil,
			queryName:       "Test",
			wantExists:      false,
		},
		{
			name:            "returns true after create",
			programToCreate: &domain.Program{Name: "Test"},
			queryName:       "Test",
			wantExists:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewProgramRepository()
			ctx := context.Background()

			if tt.programToCreate != nil {
				require.NoError(t, store.Create(ctx, tt.programToCreate))
			}

			exists, err := store.ExistsByName(ctx, tt.queryName)
			require.NoError(t, err)
			assert.Equal(t, tt.wantExists, exists)
		})
	}
}

func TestProgramMemoryStore_List(t *testing.T) {
	store := NewProgramRepository()
	ctx := context.Background()

	// Create some programs
	for _, name := range []string{"A", "B", "C"} {
		require.NoError(t, store.Create(ctx, &domain.Program{Name: name}))
	}

	t.Run("list all", func(t *testing.T) {
		programs, _, _, err := store.List(ctx, 100, "")
		require.NoError(t, err)
		assert.Len(t, programs, 3)
	})

	t.Run("list with limit", func(t *testing.T) {
		programs, cursor, hasMore, err := store.List(ctx, 2, "")
		require.NoError(t, err)
		assert.Len(t, programs, 2)
		assert.True(t, hasMore)
		assert.NotEmpty(t, cursor)
	})
}
