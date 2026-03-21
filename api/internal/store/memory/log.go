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

// logStore implements LogRepository with in-memory map storage
type logStore struct {
	mu   sync.RWMutex
	logs map[uuid.UUID]*domain.Log
}

// NewLogRepository creates a new in-memory Log repository
func NewLogRepository() repository.LogRepository {
	return &logStore{
		logs: make(map[uuid.UUID]*domain.Log),
	}
}

func (s *logStore) Create(ctx context.Context, log *domain.Log) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	log.ID = uuid.New()
	log.CreatedAt = now
	log.UpdatedAt = now

	for i := range log.Entries {
		log.Entries[i].ID = uuid.New()
		log.Entries[i].LogID = log.ID
	}

	s.logs[log.ID] = log
	return nil
}

func (s *logStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.Log, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	log, exists := s.logs[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return s.copyLog(log), nil
}

func (s *logStore) copyLog(l *domain.Log) *domain.Log {
	cp := *l
	cp.Entries = make([]domain.LogEntry, len(l.Entries))
	for i, e := range l.Entries {
		cp.Entries[i] = e
		if e.Metadata != nil {
			cp.Entries[i].Metadata = make([]byte, len(e.Metadata))
			copy(cp.Entries[i].Metadata, e.Metadata)
		}
	}
	if l.Metadata != nil {
		cp.Metadata = make([]byte, len(l.Metadata))
		copy(cp.Metadata, l.Metadata)
	}
	return &cp
}

func (s *logStore) Update(ctx context.Context, log *domain.Log) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.logs[log.ID]; !exists {
		return domain.ErrNotFound
	}

	log.UpdatedAt = time.Now()

	for i := range log.Entries {
		if log.Entries[i].ID == uuid.Nil {
			log.Entries[i].ID = uuid.New()
		}
		log.Entries[i].LogID = log.ID
	}

	s.logs[log.ID] = log
	return nil
}

func (s *logStore) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.logs[id]; !exists {
		return domain.ErrNotFound
	}

	delete(s.logs, id)
	return nil
}

func (s *logStore) List(ctx context.Context, programID *uuid.UUID, limit int, after string) ([]*domain.Log, string, bool, error) {
	return s.listWithFilter(ctx, programID, nil, nil, limit, after)
}

func (s *logStore) ListByPerformedAtRange(ctx context.Context, programID *uuid.UUID, performedAtFrom, performedAtTo *time.Time, limit int, after string) ([]*domain.Log, string, bool, error) {
	return s.listWithFilter(ctx, programID, performedAtFrom, performedAtTo, limit, after)
}

func (s *logStore) listWithFilter(ctx context.Context, programID *uuid.UUID, performedAtFrom, performedAtTo *time.Time, limit int, after string) ([]*domain.Log, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	logs := make([]*domain.Log, 0, len(s.logs))
	for _, l := range s.logs {
		if programID != nil && (l.ProgramID == nil || *l.ProgramID != *programID) {
			continue
		}
		if performedAtFrom != nil && l.PerformedAt.Before(*performedAtFrom) {
			continue
		}
		if performedAtTo != nil && !l.PerformedAt.Before(*performedAtTo) {
			continue
		}
		logs = append(logs, l)
	}

	sort.Slice(logs, func(i, j int) bool {
		if !logs[i].PerformedAt.Equal(logs[j].PerformedAt) {
			return logs[i].PerformedAt.After(logs[j].PerformedAt)
		}
		return logs[i].ID.String() < logs[j].ID.String()
	})

	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		for i, l := range logs {
			if l.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := startIdx + limit
	if endIdx > len(logs) {
		endIdx = len(logs)
	}

	result := logs[startIdx:endIdx]
	hasMore := endIdx < len(logs)

	var nextCursor string
	if hasMore && len(result) > 0 {
		nextCursor = encodeCursor(result[len(result)-1].ID)
	}

	copies := make([]*domain.Log, len(result))
	for i, l := range result {
		copies[i] = s.copyLog(l)
	}

	return copies, nextCursor, hasMore, nil
}

func (s *logStore) ListDistinctLoggedSessionsByProgramID(ctx context.Context, programID uuid.UUID) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	seen := make(map[string]struct{})
	for _, l := range s.logs {
		if l.ProgramID != nil && *l.ProgramID == programID && l.SessionName != nil {
			seen[*l.SessionName] = struct{}{}
		}
	}

	result := make([]string, 0, len(seen))
	for name := range seen {
		result = append(result, name)
	}
	return result, nil
}
