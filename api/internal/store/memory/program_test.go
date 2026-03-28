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
			programToCreate: &domain.Program{Name: "Test", Status: domain.ProgramStatusCreated},
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
