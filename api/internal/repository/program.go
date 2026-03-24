package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// ProgramRepository defines the interface for Program storage operations.
// Programs are immutable after creation (no update endpoint).
type ProgramRepository interface {
	// Create stores a new Program (with all sessions and entries) and returns it with generated IDs and timestamps
	Create(ctx context.Context, program *domain.Program) error

	// GetByID retrieves a Program by its ID including all sessions and entries, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error)

	// UpdateStatus updates the status of a Program (e.g., active → completed)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProgramStatus) error

	// Delete removes a Program by ID (cascades to sessions and entries), returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves Programs with pagination.
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// programTemplateID: optional filter by source template
	// status: optional filter by status ("active", "completed", or "" for all)
	// Returns: programs, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string, programTemplateID *uuid.UUID, status string) ([]*domain.Program, string, bool, error)

	// ExistsByProgramTemplateID checks if any Programs reference the given program template
	ExistsByProgramTemplateID(ctx context.Context, programTemplateID uuid.UUID) (bool, error)

	// ExistsByName checks if a Program with the given name already exists.
	ExistsByName(ctx context.Context, name string) (bool, error)
}
