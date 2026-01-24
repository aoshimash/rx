package memory

import (
	"context"
	"encoding/base64"
	"sort"
	"sync"
	"time"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/aoshimash/optel-training/api/internal/repository"
	"github.com/google/uuid"
)

// exerciseStore implements ExerciseRepository with in-memory map storage
type exerciseStore struct {
	mu        sync.RWMutex
	exercises map[uuid.UUID]*domain.Exercise
}

// NewExerciseRepository creates a new in-memory Exercise repository
func NewExerciseRepository() repository.ExerciseRepository {
	return &exerciseStore{
		exercises: make(map[uuid.UUID]*domain.Exercise),
	}
}

func (s *exerciseStore) Create(ctx context.Context, exercise *domain.Exercise) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	exercise.ID = uuid.New()
	exercise.CreatedAt = now
	exercise.UpdatedAt = now

	s.exercises[exercise.ID] = exercise
	return nil
}

func (s *exerciseStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Exercise, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	exercise, exists := s.exercises[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	// Return a copy to prevent external modifications
	result := *exercise
	return &result, nil
}

func (s *exerciseStore) Update(ctx context.Context, exercise *domain.Exercise) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.exercises[exercise.ID]; !exists {
		return domain.ErrNotFound
	}

	exercise.UpdatedAt = time.Now()
	s.exercises[exercise.ID] = exercise
	return nil
}

func (s *exerciseStore) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.exercises[id]; !exists {
		return domain.ErrNotFound
	}

	delete(s.exercises, id)
	return nil
}

func (s *exerciseStore) List(ctx context.Context, limit int, after string) ([]*domain.Exercise, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert map to slice and sort by ID (for consistent pagination)
	exercises := make([]*domain.Exercise, 0, len(s.exercises))
	for _, ex := range s.exercises {
		exercises = append(exercises, ex)
	}

	// Sort by ID for consistent ordering
	sort.Slice(exercises, func(i, j int) bool {
		return exercises[i].ID.String() < exercises[j].ID.String()
	})

	// Decode cursor if provided
	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		// Find the index after the cursor
		for i, ex := range exercises {
			if ex.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	// Apply limit
	endIdx := startIdx + limit
	if endIdx > len(exercises) {
		endIdx = len(exercises)
	}

	result := exercises[startIdx:endIdx]
	hasMore := endIdx < len(exercises)

	var nextCursor string
	if hasMore && len(result) > 0 {
		nextCursor = encodeCursor(result[len(result)-1].ID)
	}

	// Return copies to prevent external modifications
	copies := make([]*domain.Exercise, len(result))
	for i, ex := range result {
		copy := *ex
		copies[i] = &copy
	}

	return copies, nextCursor, hasMore, nil
}

// encodeCursor encodes a UUID as a base64 string for pagination
func encodeCursor(id uuid.UUID) string {
	return base64.URLEncoding.EncodeToString(id[:])
}

// decodeCursor decodes a base64 string to a UUID
func decodeCursor(cursor string) (uuid.UUID, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.FromBytes(data)
}
