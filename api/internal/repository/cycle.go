package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// CycleRepository defines the interface for Cycle storage operations
type CycleRepository interface {
	// Create stores a new Cycle and returns it with generated ID and timestamp
	Create(ctx context.Context, cycle *domain.Cycle) error

	// GetByID retrieves a Cycle by its ID, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Cycle, error)

	// Delete removes a Cycle by ID, returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves Cycles with pagination
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (base64-encoded UUID string)
	// Returns: cycles, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string) ([]*domain.Cycle, string, bool, error)

	// ListByProgramID retrieves Cycles filtered by program_id with pagination
	ListByProgramID(ctx context.Context, programID uuid.UUID, limit int, after string) ([]*domain.Cycle, string, bool, error)

	// ExistsByProgramID checks if any Cycles reference the given program
	ExistsByProgramID(ctx context.Context, programID uuid.UUID) (bool, error)
}
