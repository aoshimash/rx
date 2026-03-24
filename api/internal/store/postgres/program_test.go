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

	t.Run("returns false when not found", func(t *testing.T) {
		exists, err := repo.ExistsByName(ctx, "NonExistentProgram_"+t.Name())
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("returns true after create", func(t *testing.T) {
		name := "TestProgram_ExistsByName_" + t.Name()
		p := &domain.Program{
			Name:   name,
			Status: domain.ProgramStatusCreated,
		}
		require.NoError(t, repo.Create(ctx, p))
		t.Cleanup(func() {
			_ = repo.Delete(ctx, p.ID)
		})

		exists, err := repo.ExistsByName(ctx, name)
		require.NoError(t, err)
		assert.True(t, exists)
	})
}

func init() {
	if os.Getenv("DB_HOST") == "" {
		os.Setenv("DB_HOST", "localhost")
	}
}
