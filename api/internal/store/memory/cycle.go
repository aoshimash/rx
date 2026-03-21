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

// cycleStore implements CycleRepository with in-memory map storage
type cycleStore struct {
	mu     sync.RWMutex
	cycles map[uuid.UUID]*domain.Cycle
}

// NewCycleRepository creates a new in-memory Cycle repository
func NewCycleRepository() repository.CycleRepository {
	return &cycleStore{
		cycles: make(map[uuid.UUID]*domain.Cycle),
	}
}

func (s *cycleStore) Create(ctx context.Context, cycle *domain.Cycle) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cycle.ID = uuid.New()
	cycle.CreatedAt = time.Now()

	s.cycles[cycle.ID] = cycle
	return nil
}

func (s *cycleStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Cycle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cycle, exists := s.cycles[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return s.copyCycle(cycle), nil
}

func (s *cycleStore) copyCycle(c *domain.Cycle) *domain.Cycle {
	cp := *c
	if c.Metadata != nil {
		cp.Metadata = make([]byte, len(c.Metadata))
		copy(cp.Metadata, c.Metadata)
	}
	return &cp
}

func (s *cycleStore) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.cycles[id]; !exists {
		return domain.ErrNotFound
	}

	delete(s.cycles, id)
	return nil
}

func (s *cycleStore) List(ctx context.Context, limit int, after string) ([]*domain.Cycle, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cycles := make([]*domain.Cycle, 0, len(s.cycles))
	for _, c := range s.cycles {
		cycles = append(cycles, c)
	}

	return s.paginateCycles(cycles, limit, after)
}

func (s *cycleStore) ListByProgramID(ctx context.Context, programID uuid.UUID, limit int, after string) ([]*domain.Cycle, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cycles := make([]*domain.Cycle, 0)
	for _, c := range s.cycles {
		if c.ProgramID == programID {
			cycles = append(cycles, c)
		}
	}

	return s.paginateCycles(cycles, limit, after)
}

func (s *cycleStore) ExistsByProgramID(ctx context.Context, programID uuid.UUID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.cycles {
		if c.ProgramID == programID {
			return true, nil
		}
	}
	return false, nil
}

func (s *cycleStore) paginateCycles(cycles []*domain.Cycle, limit int, after string) ([]*domain.Cycle, string, bool, error) {
	sort.Slice(cycles, func(i, j int) bool {
		return cycles[i].ID.String() < cycles[j].ID.String()
	})

	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		for i, c := range cycles {
			if c.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := startIdx + limit
	if endIdx > len(cycles) {
		endIdx = len(cycles)
	}

	result := cycles[startIdx:endIdx]
	hasMore := endIdx < len(cycles)

	var nextCursor string
	if hasMore && len(result) > 0 {
		nextCursor = encodeCursor(result[len(result)-1].ID)
	}

	copies := make([]*domain.Cycle, len(result))
	for i, c := range result {
		copies[i] = s.copyCycle(c)
	}

	return copies, nextCursor, hasMore, nil
}
