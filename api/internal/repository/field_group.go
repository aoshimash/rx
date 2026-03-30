package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// FieldGroupRepository defines the interface for FieldGroup storage operations.
type FieldGroupRepository interface {
	// List returns all field groups for the given user
	List(ctx context.Context, userID string) ([]domain.FieldGroup, error)

	// GetByID retrieves a field group by ID, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.FieldGroup, error)

	// Create stores a new field group
	Create(ctx context.Context, userID string, fg *domain.FieldGroup) error

	// Update replaces the content of an existing field group
	Update(ctx context.Context, fg *domain.FieldGroup) error

	// Delete removes a field group by ID, returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error
}
