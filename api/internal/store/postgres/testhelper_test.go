package postgres

import (
	"context"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	postgresmodule "github.com/testcontainers/testcontainers-go/modules/postgres"
)

// setupTestDB creates a test PostgreSQL container and returns a connection pool
func setupTestDB(ctx context.Context) (*pgxpool.Pool, func(), error) {
	// Start PostgreSQL container
	postgresContainer, err := postgresmodule.Run(ctx,
		"postgres:17",
		postgresmodule.WithDatabase("test_db"),
		postgresmodule.WithUsername("test_user"),
		postgresmodule.WithPassword("test_password"),
		postgresmodule.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, err
	}

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, err
	}

	// Create connection pool
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, nil, err
	}

	// Run migrations
	if err := runMigrationsForTest(ctx, pool, connStr); err != nil {
		pool.Close()
		return nil, nil, err
	}

	// Cleanup function
	cleanup := func() {
		pool.Close()
		_ = postgresContainer.Terminate(ctx)
	}

	return pool, cleanup, nil
}

// runMigrationsForTest runs migrations against the test database
func runMigrationsForTest(ctx context.Context, pool *pgxpool.Pool, connStr string) error {
	// Get underlying *sql.DB for golang-migrate
	sqlDB := stdlibDB(pool)
	defer func() { _ = sqlDB.Close() }()

	// Create postgres driver instance
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}

	// Get absolute path to migrations directory
	// We need to go up from internal/store/postgres to api/ root
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	// Assuming we're in api/internal/store/postgres, go up 3 levels to api/
	apiRoot := filepath.Join(wd, "../../..")
	migrationsPath := filepath.Join(apiRoot, "migrations")

	// Create migrate instance using file:// source
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
