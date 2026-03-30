package memory

import (
	"context"
	"sync"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

type fieldGroupStore struct {
	mu     sync.RWMutex
	groups map[uuid.UUID]*domain.FieldGroup
	// userIndex maps userID to a set of field group IDs
	userIndex map[string]map[uuid.UUID]struct{}
}

func NewFieldGroupRepository() repository.FieldGroupRepository {
	return &fieldGroupStore{
		groups:    make(map[uuid.UUID]*domain.FieldGroup),
		userIndex: make(map[string]map[uuid.UUID]struct{}),
	}
}

func (s *fieldGroupStore) List(ctx context.Context, userID string) ([]domain.FieldGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids, ok := s.userIndex[userID]
	if !ok {
		return []domain.FieldGroup{}, nil
	}

	result := make([]domain.FieldGroup, 0, len(ids))
	for id := range ids {
		if fg, exists := s.groups[id]; exists {
			result = append(result, s.copyFieldGroup(fg))
		}
	}
	return result, nil
}

func (s *fieldGroupStore) GetByID(ctx context.Context, userID string, id uuid.UUID) (*domain.FieldGroup, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids, ok := s.userIndex[userID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	if _, owns := ids[id]; !owns {
		return nil, domain.ErrNotFound
	}

	fg, exists := s.groups[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	cp := s.copyFieldGroup(fg)
	return &cp, nil
}

func (s *fieldGroupStore) Create(ctx context.Context, userID string, fg *domain.FieldGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate name within user
	if ids, ok := s.userIndex[userID]; ok {
		for id := range ids {
			if existing, exists := s.groups[id]; exists && existing.Name == fg.Name {
				return &domain.DomainError{
					Code:    domain.ErrorCodeConflict,
					Message: "field group with this name already exists",
				}
			}
		}
	}

	now := time.Now()
	fg.ID = uuid.New()
	fg.CreatedAt = now
	fg.UpdatedAt = now

	cp := s.copyFieldGroup(fg)
	s.groups[fg.ID] = &cp

	if s.userIndex[userID] == nil {
		s.userIndex[userID] = make(map[uuid.UUID]struct{})
	}
	s.userIndex[userID][fg.ID] = struct{}{}

	return nil
}

func (s *fieldGroupStore) Update(ctx context.Context, userID string, fg *domain.FieldGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ids, ok := s.userIndex[userID]
	if !ok {
		return domain.ErrNotFound
	}
	if _, owns := ids[fg.ID]; !owns {
		return domain.ErrNotFound
	}

	existing, exists := s.groups[fg.ID]
	if !exists {
		return domain.ErrNotFound
	}

	// Check for duplicate name within user (excluding the current field group)
	for id := range ids {
		if id != fg.ID {
			if other, exists := s.groups[id]; exists && other.Name == fg.Name {
				return &domain.DomainError{
					Code:    domain.ErrorCodeConflict,
					Message: "field group with this name already exists",
				}
			}
		}
	}

	fg.CreatedAt = existing.CreatedAt
	fg.UpdatedAt = time.Now()

	cp := s.copyFieldGroup(fg)
	s.groups[fg.ID] = &cp
	return nil
}

func (s *fieldGroupStore) Delete(ctx context.Context, userID string, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	ids, ok := s.userIndex[userID]
	if !ok {
		return domain.ErrNotFound
	}
	if _, owns := ids[id]; !owns {
		return domain.ErrNotFound
	}

	delete(s.groups, id)
	delete(ids, id)
	if len(ids) == 0 {
		delete(s.userIndex, userID)
	}
	return nil
}

func (s *fieldGroupStore) copyFieldGroup(fg *domain.FieldGroup) domain.FieldGroup {
	cp := *fg
	if fg.Description != nil {
		d := *fg.Description
		cp.Description = &d
	}
	if fg.ProgramFields != nil {
		cp.ProgramFields = make([]domain.FieldDef, len(fg.ProgramFields))
		copy(cp.ProgramFields, fg.ProgramFields)
	}
	if fg.LogFields != nil {
		cp.LogFields = make([]domain.FieldDef, len(fg.LogFields))
		copy(cp.LogFields, fg.LogFields)
	}
	return cp
}
