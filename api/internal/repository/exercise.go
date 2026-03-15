package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// ExerciseRepository defines the interface for Exercise storage operations
type ExerciseRepository interface {
	// Create stores a new Exercise and returns it with generated ID and timestamps
	Create(ctx context.Context, exercise *domain.Exercise) error

	// GetByID retrieves an Exercise by its ID, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Exercise, error)

	// Update replaces an existing Exercise, returns domain.ErrNotFound if not found
	Update(ctx context.Context, exercise *domain.Exercise) error

	// Delete removes an Exercise by ID, returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves Exercises with pagination
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: exercises, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string) ([]*domain.Exercise, string, bool, error)
}
