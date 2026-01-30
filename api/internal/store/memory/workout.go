package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/aoshimash/optel-workout/api/internal/repository"
	"github.com/google/uuid"
)

// workoutStore implements WorkoutRepository with in-memory map storage
type workoutStore struct {
	mu       sync.RWMutex
	workouts map[uuid.UUID]*domain.Workout
}

// NewWorkoutRepository creates a new in-memory Workout repository
func NewWorkoutRepository() repository.WorkoutRepository {
	return &workoutStore{
		workouts: make(map[uuid.UUID]*domain.Workout),
	}
}

func (s *workoutStore) Create(ctx context.Context, workout *domain.Workout) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	workout.ID = uuid.New()
	workout.CreatedAt = now
	workout.UpdatedAt = now

	// Generate IDs for nested WorkoutEntry records
	for i := range workout.Entries {
		workout.Entries[i].ID = uuid.New()
		workout.Entries[i].WorkoutID = workout.ID
	}

	s.workouts[workout.ID] = workout
	return nil
}

func (s *workoutStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Workout, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workout, exists := s.workouts[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	// Return a copy to prevent external modifications
	result := *workout
	result.Entries = make([]domain.WorkoutEntry, len(workout.Entries))
	copy(result.Entries, workout.Entries)
	return &result, nil
}

func (s *workoutStore) Update(ctx context.Context, workout *domain.Workout) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.workouts[workout.ID]; !exists {
		return domain.ErrNotFound
	}

	workout.UpdatedAt = time.Now()

	// Generate IDs for new WorkoutEntry records
	for i := range workout.Entries {
		if workout.Entries[i].ID == uuid.Nil {
			workout.Entries[i].ID = uuid.New()
		}
		workout.Entries[i].WorkoutID = workout.ID
	}

	s.workouts[workout.ID] = workout
	return nil
}

func (s *workoutStore) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.workouts[id]; !exists {
		return domain.ErrNotFound
	}

	// Cascade delete: WorkoutEntry records are deleted with the Workout
	delete(s.workouts, id)
	return nil
}

func (s *workoutStore) List(ctx context.Context, limit int, after string) ([]*domain.Workout, string, bool, error) {
	return s.listWithFilter(ctx, nil, nil, limit, after)
}

func (s *workoutStore) ListByTimestampRange(ctx context.Context, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.Workout, string, bool, error) {
	return s.listWithFilter(ctx, timestampFrom, timestampTo, limit, after)
}

func (s *workoutStore) listWithFilter(ctx context.Context, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.Workout, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert map to slice and filter by timestamp if provided
	workouts := make([]*domain.Workout, 0, len(s.workouts))
	for _, w := range s.workouts {
		if timestampFrom != nil && w.Timestamp.Before(*timestampFrom) {
			continue
		}
		if timestampTo != nil && !w.Timestamp.Before(*timestampTo) {
			continue
		}
		workouts = append(workouts, w)
	}

	// Sort by timestamp (descending), then by ID for consistent ordering
	sort.Slice(workouts, func(i, j int) bool {
		if !workouts[i].Timestamp.Equal(workouts[j].Timestamp) {
			return workouts[i].Timestamp.After(workouts[j].Timestamp)
		}
		return workouts[i].ID.String() < workouts[j].ID.String()
	})

	// Decode cursor if provided
	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		// Find the index after the cursor
		for i, w := range workouts {
			if w.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	// Apply limit
	endIdx := startIdx + limit
	if endIdx > len(workouts) {
		endIdx = len(workouts)
	}

	result := workouts[startIdx:endIdx]
	hasMore := endIdx < len(workouts)

	var nextCursor string
	if hasMore && len(result) > 0 {
		nextCursor = encodeCursor(result[len(result)-1].ID)
	}

	// Return copies to prevent external modifications
	copies := make([]*domain.Workout, len(result))
	for i, w := range result {
		wCopy := *w
		wCopy.Entries = make([]domain.WorkoutEntry, len(w.Entries))
		copy(wCopy.Entries, w.Entries)
		copies[i] = &wCopy
	}

	return copies, nextCursor, hasMore, nil
}

func (s *workoutStore) ListByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]*domain.Workout, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*domain.Workout
	for _, w := range s.workouts {
		for _, entry := range w.Entries {
			if entry.ExerciseID == exerciseID {
				// Return a copy
				wCopy := *w
				wCopy.Entries = make([]domain.WorkoutEntry, len(w.Entries))
				copy(wCopy.Entries, w.Entries)
				result = append(result, &wCopy)
				break // Found at least one entry, no need to check more
			}
		}
	}

	return result, nil
}

func (s *workoutStore) ListByProgramNodeID(ctx context.Context, programNodeID uuid.UUID) ([]*domain.Workout, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*domain.Workout
	for _, w := range s.workouts {
		if w.ProgramNodeID != nil && *w.ProgramNodeID == programNodeID {
			// Return a copy
			wCopy := *w
			wCopy.Entries = make([]domain.WorkoutEntry, len(w.Entries))
			copy(wCopy.Entries, w.Entries)
			result = append(result, &wCopy)
		}
	}

	return result, nil
}
