package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/aoshimash/optel-workout/api/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// RunMigrations runs all pending migrations from the migrations directory
func RunMigrations(ctx context.Context, cfg config.DatabaseConfig) error {
	connStr := buildConnectionString(cfg)

	// Create a temporary connection for migrations
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	defer pool.Close()

	// Get underlying *sql.DB for golang-migrate
	sqlDB := stdlibDB(pool)
	defer func() { _ = sqlDB.Close() }()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Get absolute path to migrations directory
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	// Create migrate instance using file:// source
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	slog.Info("Migrations completed successfully")
	return nil
}

// GetMigrationVersion returns the current migration version
func GetMigrationVersion(ctx context.Context, cfg config.DatabaseConfig) (uint, bool, error) {
	connStr := buildConnectionString(cfg)

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create connection pool: %w", err)
	}
	defer pool.Close()

	sqlDB := stdlibDB(pool)
	defer func() { _ = sqlDB.Close() }()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migrations path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	version, dirty, err := m.Version()
	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migration version: %w", err)
	}

	return version, dirty, nil
}

// stdlibDB creates a *sql.DB from connection string for golang-migrate compatibility
func stdlibDB(pool *pgxpool.Pool) *sql.DB {
	// Get connection string from pool config
	connStr := pool.Config().ConnConfig.ConnString()

	// Open a new sql.DB connection using pgx/v5/stdlib driver
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		// This should not happen with valid connection string
		panic(fmt.Sprintf("failed to open stdlib DB: %v", err))
	}
	return db
}
