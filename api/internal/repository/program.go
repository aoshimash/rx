package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// ProgramRepository defines the interface for Program storage operations.
// Programs are immutable after creation; use Archive/Unarchive instead of Update.
type ProgramRepository interface {
	// Create stores a new Program and returns it with generated ID and timestamps
	Create(ctx context.Context, program *domain.Program) error

	// GetByID retrieves a Program by its ID including all ProgramEntry records, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error)

	// Archive sets archived_at on a Program, returns domain.ErrNotFound if not found
	Archive(ctx context.Context, id uuid.UUID) error

	// Unarchive clears archived_at on a Program, returns domain.ErrNotFound if not found
	Unarchive(ctx context.Context, id uuid.UUID) error

	// Delete removes a Program by ID (cascades to ProgramEntry records), returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves Programs with pagination.
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// includeArchived: when true, archived programs are included in results
	// Returns: programs, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string, includeArchived bool) ([]*domain.Program, string, bool, error)
}
