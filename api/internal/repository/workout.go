package repository

import (
	"context"
	"time"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/google/uuid"
)

// WorkoutRepository defines the interface for Workout storage operations
type WorkoutRepository interface {
	// Create stores a new Workout and returns it with generated ID and timestamps
	Create(ctx context.Context, workout *domain.Workout) error

	// GetByID retrieves a Workout by its ID including all WorkoutEntry records, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Workout, error)

	// Update replaces an existing Workout including all WorkoutEntry records, returns domain.ErrNotFound if not found
	Update(ctx context.Context, workout *domain.Workout) error

	// Delete removes a Workout by ID (cascades to WorkoutEntry records), returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves Workouts with pagination
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: workouts, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string) ([]*domain.Workout, string, bool, error)

	// ListByTimestampRange retrieves Workouts filtered by timestamp range with pagination
	// timestampFrom: filter workouts at or after this timestamp (inclusive)
	// timestampTo: filter workouts before this timestamp (exclusive)
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// Returns: workouts, next cursor (empty string if no more), has_more flag
	ListByTimestampRange(ctx context.Context, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.Workout, string, bool, error)

	// ListByExerciseID retrieves Workouts that contain WorkoutEntry records referencing the given Exercise ID
	// Used for referential integrity checks when deleting an Exercise
	ListByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]*domain.Workout, error)

	// ListByProgramNodeID retrieves Workouts that reference the given ProgramNode ID
	// Used for referential integrity checks when deleting a Program
	ListByProgramNodeID(ctx context.Context, programNodeID uuid.UUID) ([]*domain.Workout, error)
}
