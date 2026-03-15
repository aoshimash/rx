package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// ProgramRepository defines the interface for Program storage operations
type ProgramRepository interface {
	// Create stores a new Program and returns it with generated ID and timestamps
	Create(ctx context.Context, program *domain.Program) error

	// GetByID retrieves a Program by its ID including all ProgramEntry records, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error)

	// Update replaces an existing Program including all ProgramEntry records, returns domain.ErrNotFound if not found
	Update(ctx context.Context, program *domain.Program) error

	// Delete removes a Program by ID (cascades to ProgramEntry records), returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves Programs with pagination
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: programs, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string) ([]*domain.Program, string, bool, error)
}
