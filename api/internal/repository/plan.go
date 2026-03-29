package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
)

// PlanRepository defines the interface for Plan storage operations.
// Each user has at most one Plan.
type PlanRepository interface {
	// GetByUserID retrieves the Plan for the given user, returns domain.ErrNotFound if no plan exists
	GetByUserID(ctx context.Context, userID string) (*domain.Plan, error)

	// Create stores a new Plan for the given user, returns conflict error if a plan already exists
	Create(ctx context.Context, userID string, plan *domain.Plan) error

	// Update replaces the content of an existing Plan, returns domain.ErrNotFound if no plan exists
	Update(ctx context.Context, userID string, plan *domain.Plan) error

	// Delete removes the Plan for the given user, returns domain.ErrNotFound if no plan exists
	Delete(ctx context.Context, userID string) error

	// AddSessions appends sessions to an existing Plan, returns domain.ErrNotFound if no plan exists
	AddSessions(ctx context.Context, userID string, sessions []domain.PlanSession) error

	// UpdateSession replaces a specific session by ID, returns domain.ErrNotFound if plan or session not found
	UpdateSession(ctx context.Context, userID string, session *domain.PlanSession) error

	// DeleteSession removes a specific session by ID, returns domain.ErrNotFound if plan or session not found
	DeleteSession(ctx context.Context, userID string, sessionID string) error
}
