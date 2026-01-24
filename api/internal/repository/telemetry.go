package repository

import (
	"context"
	"time"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/google/uuid"
)

// TelemetryPointRepository defines the interface for TelemetryPoint storage operations
type TelemetryPointRepository interface {
	// Create stores a new TelemetryPoint and returns it with generated ID and timestamps
	Create(ctx context.Context, point *domain.TelemetryPoint) error

	// GetByID retrieves a TelemetryPoint by its ID, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.TelemetryPoint, error)

	// Update replaces an existing TelemetryPoint, returns domain.ErrNotFound if not found
	Update(ctx context.Context, point *domain.TelemetryPoint) error

	// Delete removes a TelemetryPoint by ID, returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves TelemetryPoints with pagination
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: points, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string) ([]*domain.TelemetryPoint, string, bool, error)

	// ListByMetricAndTimeRange retrieves TelemetryPoints filtered by metric name and timestamp range with pagination
	// metricName: filter by metric name (required)
	// timestampFrom: filter points at or after this timestamp (inclusive)
	// timestampTo: filter points before this timestamp (exclusive)
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: points, next cursor (empty string if no more), has_more flag
	ListByMetricAndTimeRange(ctx context.Context, metricName string, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.TelemetryPoint, string, bool, error)
}
