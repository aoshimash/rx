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

// planStore implements PlanRepository with in-memory map storage
type planStore struct {
	mu    sync.RWMutex
	plans map[uuid.UUID]*domain.Plan
}

// NewPlanRepository creates a new in-memory Plan repository
func NewPlanRepository() repository.PlanRepository {
	return &planStore{
		plans: make(map[uuid.UUID]*domain.Plan),
	}
}

func (s *planStore) Create(ctx context.Context, plan *domain.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	plan.ID = uuid.New()
	plan.CreatedAt = now
	plan.UpdatedAt = now

	// Generate IDs for entries
	for i := range plan.Entries {
		plan.Entries[i].ID = uuid.New()
		plan.Entries[i].PlanID = plan.ID
	}

	s.plans[plan.ID] = plan
	return nil
}

func (s *planStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plan, exists := s.plans[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return s.copyPlan(plan), nil
}

func (s *planStore) copyPlan(p *domain.Plan) *domain.Plan {
	cp := *p
	cp.Entries = make([]domain.PlanEntry, len(p.Entries))
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

func (s *planStore) Update(ctx context.Context, plan *domain.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.plans[plan.ID]; !exists {
		return domain.ErrNotFound
	}

	plan.UpdatedAt = time.Now()

	// Generate IDs for entries
	for i := range plan.Entries {
		plan.Entries[i].ID = uuid.New()
		plan.Entries[i].PlanID = plan.ID
	}

	s.plans[plan.ID] = plan
	return nil
}

func (s *planStore) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.plans[id]; !exists {
		return domain.ErrNotFound
	}

	delete(s.plans, id)
	return nil
}

func (s *planStore) List(ctx context.Context, limit int, after string) ([]*domain.Plan, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plans := make([]*domain.Plan, 0, len(s.plans))
	for _, p := range s.plans {
		plans = append(plans, p)
	}

	sort.Slice(plans, func(i, j int) bool {
		return plans[i].ID.String() < plans[j].ID.String()
	})

	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		for i, p := range plans {
			if p.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := startIdx + limit
	if endIdx > len(plans) {
		endIdx = len(plans)
	}

	result := plans[startIdx:endIdx]
	hasMore := endIdx < len(plans)

	var nextCursor string
	if hasMore && len(result) > 0 {
		nextCursor = encodeCursor(result[len(result)-1].ID)
	}

	copies := make([]*domain.Plan, len(result))
	for i, p := range result {
		copies[i] = s.copyPlan(p)
	}

	return copies, nextCursor, hasMore, nil
}

func (s *planStore) ListByProgramID(ctx context.Context, programID uuid.UUID, limit int, after string) ([]*domain.Plan, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plans := make([]*domain.Plan, 0)
	for _, p := range s.plans {
		if p.ProgramID != nil && *p.ProgramID == programID {
			plans = append(plans, p)
		}
	}

	sort.Slice(plans, func(i, j int) bool {
		return plans[i].ID.String() < plans[j].ID.String()
	})

	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		for i, p := range plans {
			if p.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := startIdx + limit
	if endIdx > len(plans) {
		endIdx = len(plans)
	}

	result := plans[startIdx:endIdx]
	hasMore := endIdx < len(plans)

	var nextCursor string
	if hasMore && len(result) > 0 {
		nextCursor = encodeCursor(result[len(result)-1].ID)
	}

	copies := make([]*domain.Plan, len(result))
	for i, p := range result {
		copies[i] = s.copyPlan(p)
	}

	return copies, nextCursor, hasMore, nil
}
