package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// PlanRepository defines the interface for Plan storage operations
type PlanRepository interface {
	// Create stores a new Plan and returns it with generated ID and timestamps
	Create(ctx context.Context, plan *domain.Plan) error

	// GetByID retrieves a Plan by its ID including all PlanEntry records, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error)

	// Update replaces an existing Plan including all PlanEntry records, returns domain.ErrNotFound if not found
	Update(ctx context.Context, plan *domain.Plan) error

	// Delete removes a Plan by ID (cascades to PlanEntry records), returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves Plans with pagination
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: plans, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string) ([]*domain.Plan, string, bool, error)
}
