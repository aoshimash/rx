package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aoshimash/optel-workout/api/internal/config"
	"github.com/cenkalti/backoff/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	maxRetries = 5
)

// DB wraps a PostgreSQL connection pool with retry logic
type DB struct {
	pool *pgxpool.Pool
}

// NewDB creates a new database connection pool with exponential backoff retry
func NewDB(ctx context.Context, cfg config.DatabaseConfig) (*DB, error) {
	connStr := buildConnectionString(cfg)

	var pool *pgxpool.Pool
	var err error

	// Exponential backoff retry logic
	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = 1 * time.Second
	backoffConfig.MaxInterval = 30 * time.Second
	backoffConfig.MaxElapsedTime = 2 * time.Minute

	retryCount := 0
	operation := func() error {
		retryCount++
		slog.Info("Attempting database connection",
			"attempt", retryCount,
			"max_retries", maxRetries,
		)

		poolConfig, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			return fmt.Errorf("failed to parse connection string: %w", err)
		}

		// Configure connection pool
		poolConfig.MaxConns = int32(cfg.MaxConns)
		poolConfig.MinConns = int32(cfg.MinConns)
		poolConfig.MaxConnLifetime = time.Hour
		poolConfig.MaxConnIdleTime = 30 * time.Minute
		poolConfig.HealthCheckPeriod = 1 * time.Minute

		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			slog.Warn("Database connection failed",
				"attempt", retryCount,
				"error", err,
			)
			return fmt.Errorf("failed to create connection pool: %w", err)
		}

		// Test connection
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			slog.Warn("Database ping failed",
				"attempt", retryCount,
				"error", err,
			)
			return fmt.Errorf("failed to ping database: %w", err)
		}

		slog.Info("Database connection established",
			"attempt", retryCount,
			"max_conns", cfg.MaxConns,
			"min_conns", cfg.MinConns,
		)
		return nil
	}

	// Retry with exponential backoff
	err = backoff.Retry(operation, backoffConfig)
	if err != nil {
		if retryCount >= maxRetries {
			return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxRetries, err)
		}
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &DB{pool: pool}, nil
}

// Pool returns the underlying connection pool
func (db *DB) Pool() *pgxpool.Pool {
	return db.pool
}

// Ping checks database connectivity
func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Close closes all connections in the pool
func (db *DB) Close() {
	db.pool.Close()
}

// buildConnectionString constructs a PostgreSQL connection string from config
func buildConnectionString(cfg config.DatabaseConfig) string {
	// Support DATABASE_URL environment variable (for production)
	if dbURL := getEnvOrDefault("DATABASE_URL", ""); dbURL != "" {
		return dbURL
	}

	// Build connection string from individual components
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
	)
}

// getEnvOrDefault is a helper to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
