package memory

import (
	"context"
	"sync"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
)

// planStore implements PlanRepository with in-memory map storage
type planStore struct {
	mu    sync.RWMutex
	plans map[string]*domain.Plan // keyed by userID
}

// NewPlanRepository creates a new in-memory Plan repository
func NewPlanRepository() repository.PlanRepository {
	return &planStore{
		plans: make(map[string]*domain.Plan),
	}
}

func (s *planStore) Create(ctx context.Context, userID string, plan *domain.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.plans[userID]; exists {
		return &domain.DomainError{
			Code:    domain.ErrorCodeConflict,
			Message: "plan already exists for user",
		}
	}

	now := time.Now()
	plan.ID = uuid.New()
	plan.CreatedAt = now
	plan.UpdatedAt = now

	for i := range plan.Sessions {
		plan.Sessions[i].ID = uuid.New()
		plan.Sessions[i].PlanID = plan.ID
		for j := range plan.Sessions[i].Entries {
			plan.Sessions[i].Entries[j].ID = uuid.New()
			plan.Sessions[i].Entries[j].SessionID = plan.Sessions[i].ID
		}
	}

	s.plans[userID] = s.copyPlan(plan)
	return nil
}

func (s *planStore) GetByUserID(ctx context.Context, userID string) (*domain.Plan, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	plan, exists := s.plans[userID]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return s.copyPlan(plan), nil
}

func (s *planStore) Update(ctx context.Context, userID string, plan *domain.Plan) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.plans[userID]
	if !exists {
		return domain.ErrNotFound
	}

	plan.ID = existing.ID
	plan.CreatedAt = existing.CreatedAt
	plan.UpdatedAt = time.Now()

	for i := range plan.Sessions {
		if plan.Sessions[i].ID == uuid.Nil {
			plan.Sessions[i].ID = uuid.New()
		}
		plan.Sessions[i].PlanID = plan.ID
		for j := range plan.Sessions[i].Entries {
			if plan.Sessions[i].Entries[j].ID == uuid.Nil {
				plan.Sessions[i].Entries[j].ID = uuid.New()
			}
			plan.Sessions[i].Entries[j].SessionID = plan.Sessions[i].ID
		}
	}

	s.plans[userID] = s.copyPlan(plan)
	return nil
}

func (s *planStore) Delete(ctx context.Context, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.plans[userID]; !exists {
		return domain.ErrNotFound
	}

	delete(s.plans, userID)
	return nil
}

func (s *planStore) AddSessions(ctx context.Context, userID string, sessions []domain.PlanSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plan, exists := s.plans[userID]
	if !exists {
		return domain.ErrNotFound
	}

	for i := range sessions {
		sessions[i].ID = uuid.New()
		sessions[i].PlanID = plan.ID
		for j := range sessions[i].Entries {
			sessions[i].Entries[j].ID = uuid.New()
			sessions[i].Entries[j].SessionID = sessions[i].ID
		}
	}

	plan.Sessions = append(plan.Sessions, sessions...)
	plan.UpdatedAt = time.Now()
	return nil
}

func (s *planStore) UpdateSession(ctx context.Context, userID string, session *domain.PlanSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plan, exists := s.plans[userID]
	if !exists {
		return domain.ErrNotFound
	}

	for i, sess := range plan.Sessions {
		if sess.ID == session.ID {
			session.PlanID = plan.ID
			for j := range session.Entries {
				if session.Entries[j].ID == uuid.Nil {
					session.Entries[j].ID = uuid.New()
				}
				session.Entries[j].SessionID = session.ID
			}
			plan.Sessions[i] = *session
			plan.UpdatedAt = time.Now()
			return nil
		}
	}

	return domain.ErrNotFound
}

func (s *planStore) DeleteSession(ctx context.Context, userID string, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	plan, exists := s.plans[userID]
	if !exists {
		return domain.ErrNotFound
	}

	parsedID, err := uuid.Parse(sessionID)
	if err != nil {
		return domain.ErrNotFound
	}

	for i, sess := range plan.Sessions {
		if sess.ID == parsedID {
			plan.Sessions = append(plan.Sessions[:i], plan.Sessions[i+1:]...)
			plan.UpdatedAt = time.Now()
			return nil
		}
	}

	return domain.ErrNotFound
}

func (s *planStore) copyPlan(p *domain.Plan) *domain.Plan {
	cp := *p
	if p.Sessions != nil {
		cp.Sessions = make([]domain.PlanSession, len(p.Sessions))
		for i, sess := range p.Sessions {
			cp.Sessions[i] = sess
			if sess.SourceProgramID != nil {
				id := *sess.SourceProgramID
				cp.Sessions[i].SourceProgramID = &id
			}
			if sess.SourceSessionID != nil {
				id := *sess.SourceSessionID
				cp.Sessions[i].SourceSessionID = &id
			}
			if sess.Entries != nil {
				cp.Sessions[i].Entries = make([]domain.PlanSessionEntry, len(sess.Entries))
				for j, e := range sess.Entries {
					cp.Sessions[i].Entries[j] = e
					if e.Fields != nil {
						fields := make(map[string]interface{}, len(e.Fields))
						for k, v := range e.Fields {
							fields[k] = v
						}
						cp.Sessions[i].Entries[j].Fields = fields
					}
				}
			}
		}
	}
	if p.ProgramID != nil {
		id := *p.ProgramID
		cp.ProgramID = &id
	}
	return &cp
}
