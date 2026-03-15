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

	// Generate IDs for entries
	s.assignEntryIDs(program.ID, program.Entries)

	s.programs[program.ID] = program
	return nil
}

func (s *programStore) assignEntryIDs(programID uuid.UUID, entries []domain.ProgramEntry) {
	for i := range entries {
		entries[i].ID = uuid.New()
		entries[i].ProgramID = programID
	}
}

func (s *programStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	program, exists := s.programs[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	// Return a copy to prevent external modifications
	result := s.copyProgram(program)
	return result, nil
}

func (s *programStore) copyProgram(p *domain.Program) *domain.Program {
	cp := *p
	cp.Entries = make([]domain.ProgramEntry, len(p.Entries))
	copy(cp.Entries, p.Entries)
	return &cp
}

func (s *programStore) Update(ctx context.Context, program *domain.Program) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.programs[program.ID]; !exists {
		return domain.ErrNotFound
	}

	program.UpdatedAt = time.Now()

	// Generate IDs for entries
	s.assignEntryIDs(program.ID, program.Entries)

	s.programs[program.ID] = program
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

func (s *programStore) List(ctx context.Context, limit int, after string) ([]*domain.Program, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Convert map to slice and sort by ID (for consistent pagination)
	programs := make([]*domain.Program, 0, len(s.programs))
	for _, p := range s.programs {
		programs = append(programs, p)
	}

	// Sort by ID for consistent ordering
	sort.Slice(programs, func(i, j int) bool {
		return programs[i].ID.String() < programs[j].ID.String()
	})

	// Decode cursor if provided
	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		// Find the index after the cursor
		for i, p := range programs {
			if p.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	// Apply limit
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

	// Return copies to prevent external modifications
	copies := make([]*domain.Program, len(result))
	for i, p := range result {
		copies[i] = s.copyProgram(p)
	}

	return copies, nextCursor, hasMore, nil
}
