package memory

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// programStore implements ProgramRepository with in-memory map storage
type programStore struct {
	mu       sync.RWMutex
	programs map[uuid.UUID]*domain.Program
}

// NewProgramRepository creates a new in-memory Program repository
func NewProgramRepository() repository.ProgramRepository {
	return &programStore{
		programs: make(map[uuid.UUID]*domain.Program),
	}
}

func (s *programStore) Create(ctx context.Context, program *domain.Program) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	program.ID = uuid.New()
	program.CreatedAt = now
	program.UpdatedAt = now

	for i := range program.Entries {
		program.Entries[i].ID = uuid.New()
		program.Entries[i].ProgramID = program.ID
	}

	s.programs[program.ID] = program
	return nil
}

func (s *programStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	program, exists := s.programs[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return s.copyProgram(program), nil
}

func (s *programStore) copyProgram(p *domain.Program) *domain.Program {
	cp := *p
	cp.Entries = make([]domain.ProgramEntry, len(p.Entries))
	for i, e := range p.Entries {
		cp.Entries[i] = e
		if e.Metadata != nil {
			cp.Entries[i].Metadata = make([]byte, len(e.Metadata))
			copy(cp.Entries[i].Metadata, e.Metadata)
		}
	}
	if p.Metadata != nil {
		cp.Metadata = make([]byte, len(p.Metadata))
		copy(cp.Metadata, p.Metadata)
	}
	return &cp
}

func (s *programStore) Archive(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	program, exists := s.programs[id]
	if !exists {
		return domain.ErrNotFound
	}

	now := time.Now()
	program.ArchivedAt = &now
	return nil
}

func (s *programStore) Unarchive(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	program, exists := s.programs[id]
	if !exists {
		return domain.ErrNotFound
	}

	program.ArchivedAt = nil
	return nil
}

func (s *programStore) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.programs[id]; !exists {
		return domain.ErrNotFound
	}

	delete(s.programs, id)
	return nil
}

func (s *programStore) List(ctx context.Context, limit int, after string, includeArchived bool) ([]*domain.Program, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	programs := make([]*domain.Program, 0, len(s.programs))
	for _, p := range s.programs {
		if !includeArchived && p.ArchivedAt != nil {
			continue
		}
		programs = append(programs, p)
	}

	sort.Slice(programs, func(i, j int) bool {
		return programs[i].ID.String() < programs[j].ID.String()
	})

	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		for i, p := range programs {
			if p.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := startIdx + limit
	if endIdx > len(programs) {
		endIdx = len(programs)
	}

	result := programs[startIdx:endIdx]
	hasMore := endIdx < len(programs)

	var nextCursor string
	if hasMore && len(result) > 0 {
		nextCursor = encodeCursor(result[len(result)-1].ID)
	}

	copies := make([]*domain.Program, len(result))
	for i, p := range result {
		copies[i] = s.copyProgram(p)
	}

	return copies, nextCursor, hasMore, nil
}
