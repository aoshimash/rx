package repository

import (
	"context"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// LogRepository defines the interface for Log storage operations
type LogRepository interface {
	// Create stores a new Log and returns it with generated ID and timestamps
	Create(ctx context.Context, log *domain.Log) error

	// GetByID retrieves a Log by its ID including all LogEntry records, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Log, error)

	// Update replaces an existing Log including all LogEntry records, returns domain.ErrNotFound if not found
	Update(ctx context.Context, log *domain.Log) error

	// Delete removes a Log by ID (cascades to LogEntry records), returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves Logs with pagination
	// programID: optional filter by program (nil means no filter)
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: logs, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, programID *uuid.UUID, limit int, after string) ([]*domain.Log, string, bool, error)

	// ListByPerformedAtRange retrieves Logs filtered by performed_at range with pagination
	// programID: optional filter by program (nil means no filter)
	// performedAtFrom: filter logs at or after this timestamp (inclusive)
	// performedAtTo: filter logs before this timestamp (exclusive)
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: logs, next cursor (empty string if no more), has_more flag
	ListByPerformedAtRange(ctx context.Context, programID *uuid.UUID, performedAtFrom, performedAtTo *time.Time, limit int, after string) ([]*domain.Log, string, bool, error)

	// ListDistinctLoggedSessionsByProgramID returns distinct session_names that have at least one log for the given program.
	// Used to check program completion when a new log is created.
	ListDistinctLoggedSessionsByProgramID(ctx context.Context, programID uuid.UUID) ([]string, error)
}
