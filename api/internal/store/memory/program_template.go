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

// programTemplateStore implements ProgramTemplateRepository with in-memory map storage
type programTemplateStore struct {
	mu        sync.RWMutex
	templates map[uuid.UUID]*domain.ProgramTemplate
}

// NewProgramTemplateRepository creates a new in-memory ProgramTemplate repository
func NewProgramTemplateRepository() repository.ProgramTemplateRepository {
	return &programTemplateStore{
		templates: make(map[uuid.UUID]*domain.ProgramTemplate),
	}
}

func (s *programTemplateStore) Create(ctx context.Context, tmpl *domain.ProgramTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	tmpl.ID = uuid.New()
	tmpl.CreatedAt = now
	tmpl.UpdatedAt = now

	for i := range tmpl.Entries {
		tmpl.Entries[i].ID = uuid.New()
		tmpl.Entries[i].ProgramTemplateID = tmpl.ID
	}

	s.templates[tmpl.ID] = tmpl
	return nil
}

func (s *programTemplateStore) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProgramTemplate, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tmpl, exists := s.templates[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	return s.copyTemplate(tmpl), nil
}

func (s *programTemplateStore) copyTemplate(t *domain.ProgramTemplate) *domain.ProgramTemplate {
	cp := *t
	cp.Entries = make([]domain.ProgramTemplateEntry, len(t.Entries))
	for i, e := range t.Entries {
		cp.Entries[i] = e
		if e.Metadata != nil {
			cp.Entries[i].Metadata = make([]byte, len(e.Metadata))
			copy(cp.Entries[i].Metadata, e.Metadata)
		}
	}
	if t.Metadata != nil {
		cp.Metadata = make([]byte, len(t.Metadata))
		copy(cp.Metadata, t.Metadata)
	}
	return &cp
}

func (s *programTemplateStore) CreateAndArchive(ctx context.Context, tmpl *domain.ProgramTemplate, archiveID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.templates[archiveID]
	if !exists {
		return domain.ErrNotFound
	}

	now := time.Now()
	tmpl.ID = uuid.New()
	tmpl.CreatedAt = now
	tmpl.UpdatedAt = now

	for i := range tmpl.Entries {
		tmpl.Entries[i].ID = uuid.New()
		tmpl.Entries[i].ProgramTemplateID = tmpl.ID
	}

	s.templates[tmpl.ID] = tmpl

	archivedAt := now
	existing.ArchivedAt = &archivedAt

	return nil
}

func (s *programTemplateStore) Archive(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tmpl, exists := s.templates[id]
	if !exists {
		return domain.ErrNotFound
	}

	now := time.Now()
	tmpl.ArchivedAt = &now
	return nil
}

func (s *programTemplateStore) Unarchive(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tmpl, exists := s.templates[id]
	if !exists {
		return domain.ErrNotFound
	}

	tmpl.ArchivedAt = nil
	return nil
}

func (s *programTemplateStore) Delete(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.templates[id]; !exists {
		return domain.ErrNotFound
	}

	delete(s.templates, id)
	return nil
}

func (s *programTemplateStore) Update(ctx context.Context, tmpl *domain.ProgramTemplate) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, exists := s.templates[tmpl.ID]
	if !exists {
		return domain.ErrNotFound
	}

	now := time.Now()
	tmpl.CreatedAt = existing.CreatedAt
	tmpl.CreatedBy = existing.CreatedBy
	tmpl.ArchivedAt = existing.ArchivedAt
	tmpl.SourceTemplateID = existing.SourceTemplateID
	tmpl.UpdatedAt = now

	for i := range tmpl.Entries {
		tmpl.Entries[i].ID = uuid.New()
		tmpl.Entries[i].ProgramTemplateID = tmpl.ID
	}

	s.templates[tmpl.ID] = tmpl
	return nil
}

func (s *programTemplateStore) List(ctx context.Context, limit int, after string, includeArchived bool) ([]*domain.ProgramTemplate, string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	templates := make([]*domain.ProgramTemplate, 0, len(s.templates))
	for _, t := range s.templates {
		if !includeArchived && t.ArchivedAt != nil {
			continue
		}
		templates = append(templates, t)
	}

	sort.Slice(templates, func(i, j int) bool {
		return templates[i].ID.String() < templates[j].ID.String()
	})

	var startIdx int
	if after != "" {
		cursorID, err := decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		for i, t := range templates {
			if t.ID == cursorID {
				startIdx = i + 1
				break
			}
		}
	}

	endIdx := startIdx + limit
	if endIdx > len(templates) {
		endIdx = len(templates)
	}

	result := templates[startIdx:endIdx]
	hasMore := endIdx < len(templates)

	var nextCursor string
	if hasMore && len(result) > 0 {
		nextCursor = encodeCursor(result[len(result)-1].ID)
	}

	copies := make([]*domain.ProgramTemplate, len(result))
	for i, t := range result {
		copies[i] = s.copyTemplate(t)
	}

	return copies, nextCursor, hasMore, nil
}

func (s *programTemplateStore) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.templates[id]
	return exists, nil
}
