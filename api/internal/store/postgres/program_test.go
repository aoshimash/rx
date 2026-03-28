//go:build integration

package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/aoshimash/rx/api/internal/config"
	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestDB(t *testing.T) *DB {
	t.Helper()
	cfg := config.DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		Name:     getEnvOrDefault("DB_NAME", "rx_test"),
		SSLMode:  "disable",
		MaxConns: 5,
		MinConns: 1,
	}
	db, err := NewDB(context.Background(), cfg)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	return db
}

func TestProgramRepository_ExistsByName(t *testing.T) {
	db := newTestDB(t)
	repo := NewProgramRepository(db.Pool())
	ctx := context.Background()

	existing := &domain.Program{
		Name:   "ExistingProgram_ExistsByName",
		Status: domain.ProgramStatusCreated,
	}
	require.NoError(t, repo.Create(ctx, existing))
	t.Cleanup(func() { _ = repo.Delete(ctx, existing.ID) })

	tests := []struct {
		name       string
		searchName string
		wantExists bool
	}{
		{
			name:       "returns false when program does not exist",
			searchName: "NonExistentProgram_ExistsByName",
			wantExists: false,
		},
		{
			name:       "returns true when program exists",
			searchName: existing.Name,
			wantExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists, err := repo.ExistsByName(ctx, tt.searchName)
			require.NoError(t, err)
			assert.Equal(t, tt.wantExists, exists)
		})
	}
}

func init() {
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}
}
