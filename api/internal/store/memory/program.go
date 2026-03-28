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

	for i := range program.Sessions {
		program.Sessions[i].ID = uuid.New()
		program.Sessions[i].ProgramID = program.ID
		for j := range program.Sessions[i].Entries {
			program.Sessions[i].Entries[j].ID = uuid.New()
			program.Sessions[i].Entries[j].SessionID = program.Sessions[i].ID
		}
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
	cp.Sessions = make([]domain.ProgramSession, len(p.Sessions))
	for i, sess := range p.Sessions {
		cp.Sessions[i] = sess
		cp.Sessions[i].Entries = make([]domain.ProgramSessionEntry, len(sess.Entries))
		for j, e := range sess.Entries {
			cp.Sessions[i].Entries[j] = e
			if e.Metadata != nil {
				cp.Sessions[i].Entries[j].Metadata = make([]byte, len(e.Metadata))
				copy(cp.Sessions[i].Entries[j].Metadata, e.Metadata)
			}
		}
	}
	if p.Metadata != nil {
		cp.Metadata = make([]byte, len(p.Metadata))
		copy(cp.Metadata, p.Metadata)
	}
	return &cp
}

func (s *programStore) Update(ctx context.Context, program *domain.Program) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.programs[program.ID]
	if !exists {
		return domain.ErrNotFound
	}

	program.CreatedAt = existing.CreatedAt
	program.UpdatedAt = time.Now()

	for i := range program.Sessions {
		program.Sessions[i].ID = uuid.New()
		program.Sessions[i].ProgramID = program.ID
		for j := range program.Sessions[i].Entries {
			program.Sessions[i].Entries[j].ID = uuid.New()
			program.Sessions[i].Entries[j].SessionID = program.Sessions[i].ID
		}
	}

	s.programs[program.ID] = program
	return nil
}

func (s *programStore) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProgramStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	program, exists := s.programs[id]
	if !exists {
		return domain.ErrNotFound
	}

	program.Status = status
	program.UpdatedAt = time.Now()
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

func (s *programStore) List(ctx context.Context, limit int, after string, status string) ([]*domain.Program, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	programs := make([]*domain.Program, 0, len(s.programs))
	for _, p := range s.programs {
		if status != "" && string(p.Status) != status {
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

func (s *programStore) ExistsByName(ctx context.Context, name string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.programs {
		if p.Name == name {
			return true, nil
		}
	}
	return false, nil
}
