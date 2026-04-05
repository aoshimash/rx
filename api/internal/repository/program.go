package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// ProgramRepository defines the interface for Program storage operations.
type ProgramRepository interface {
	// Create stores a new Program (with all sessions and entries) and returns it with generated IDs and timestamps
	Create(ctx context.Context, program *domain.Program) error

	// Update replaces the content of an existing Program (name, notes, sessions, entries, groups).
	// The program must already exist; returns domain.ErrNotFound if not found.
	// Sessions and entries are fully replaced (delete-and-reinsert).
	Update(ctx context.Context, program *domain.Program) error

	// GetByID retrieves a Program by its ID including all sessions and entries, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error)

	// Delete removes a Program by ID (cascades to sessions and entries), returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves Programs with pagination.
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: programs, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string) ([]*domain.Program, string, bool, error)

	// UpdateStatus updates only the status of an existing Program.
	// Returns domain.ErrNotFound if not found.
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProgramStatus) (*domain.Program, error)

	// ExistsByName checks if a Program with the given name already exists.
	ExistsByName(ctx context.Context, name string) (bool, error)
}
