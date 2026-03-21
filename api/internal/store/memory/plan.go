package memory

import (
	"context"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// planStore implements PlanRepository with in-memory map storage
type planStore struct {
	mu      sync.RWMutex
	plans   map[uuid.UUID]*domain.Plan
	logRepo repository.LogRepository
}

// NewPlanRepository creates a new in-memory Plan repository.
// logRepo is used to filter out executed plans (plans with associated logs) from list results.
func NewPlanRepository(logRepo repository.LogRepository) repository.PlanRepository {
	return &planStore{
		plans:   make(map[uuid.UUID]*domain.Plan),
		logRepo: logRepo,
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
	allPlans := make([]*domain.Plan, 0, len(s.plans))
	for _, p := range s.plans {
		allPlans = append(allPlans, p)
	}
	s.mu.RUnlock()

	// Filter out executed plans (plans with associated logs)
	plans := make([]*domain.Plan, 0, len(allPlans))
	for _, p := range allPlans {
		logs, err := s.logRepo.ListByPlanID(ctx, p.ID)
		if err != nil {
			return nil, "", false, err
		}
		if len(logs) == 0 {
			plans = append(plans, p)
		}
	}

	return s.paginatePlans(plans, limit, after)
}

func (s *planStore) ListByCycleID(ctx context.Context, cycleID uuid.UUID, limit int, after string) ([]*domain.Plan, string, bool, error) {
	s.mu.RLock()
	allPlans := make([]*domain.Plan, 0)
	for _, p := range s.plans {
		if p.CycleID != nil && *p.CycleID == cycleID {
			allPlans = append(allPlans, p)
		}
	}
	s.mu.RUnlock()

	// Filter out executed plans (plans with associated logs)
	plans := make([]*domain.Plan, 0, len(allPlans))
	for _, p := range allPlans {
		logs, err := s.logRepo.ListByPlanID(ctx, p.ID)
		if err != nil {
			return nil, "", false, err
		}
		if len(logs) == 0 {
			plans = append(plans, p)
		}
	}

	return s.paginatePlans(plans, limit, after)
}

// paginatePlans sorts plans by (created_at, id) and applies cursor-based pagination.
func (s *planStore) paginatePlans(plans []*domain.Plan, limit int, after string) ([]*domain.Plan, string, bool, error) {
	sort.Slice(plans, func(i, j int) bool {
		if !plans[i].CreatedAt.Equal(plans[j].CreatedAt) {
			return plans[i].CreatedAt.Before(plans[j].CreatedAt)
		}
		return plans[i].ID.String() < plans[j].ID.String()
	})

	var startIdx int
	if after != "" {
		cursorTime, cursorID, err := decodePlanCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		startIdx = len(plans) // default: cursor is past all elements
		for i, p := range plans {
			if p.CreatedAt.After(cursorTime) || (p.CreatedAt.Equal(cursorTime) && p.ID.String() > cursorID.String()) {
				startIdx = i
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
		last := result[len(result)-1]
		nextCursor = encodePlanCursor(last.CreatedAt, last.ID)
	}

	copies := make([]*domain.Plan, len(result))
	for i, p := range result {
		copies[i] = s.copyPlan(p)
	}

	return copies, nextCursor, hasMore, nil
}

func (s *planStore) CountByCycleID(ctx context.Context, cycleID uuid.UUID) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, p := range s.plans {
		if p.CycleID != nil && *p.CycleID == cycleID {
			count++
		}
	}
	return count, nil
}

// encodePlanCursor encodes a (created_at, id) pair for plan cursor-based pagination
func encodePlanCursor(createdAt time.Time, id uuid.UUID) string {
	s := createdAt.Format(time.RFC3339Nano) + "|" + id.String()
	return base64.URLEncoding.EncodeToString([]byte(s))
}

// decodePlanCursor decodes a plan cursor to (created_at, id)
func decodePlanCursor(cursor string) (time.Time, uuid.UUID, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	parts := strings.SplitN(string(data), "|", 2)
	if len(parts) != 2 {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid plan cursor format")
	}
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	id, err := uuid.Parse(parts[1])
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	return t, id, nil
}
