package repository

import (
	"context"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// ProgramTemplateRepository defines the interface for ProgramTemplate storage operations.
type ProgramTemplateRepository interface {
	// Create stores a new ProgramTemplate and returns it with generated ID and timestamps
	Create(ctx context.Context, tmpl *domain.ProgramTemplate) error

	// Update replaces the content of an existing ProgramTemplate (name, description, notes, metadata, weeks, days_per_week, entries).
	// The template ID, source_template_id, created_at, created_by, and archived_at are preserved and must not be changed by Update.
	// updated_at is set to NOW(). Entries are replaced entirely (old entries deleted, new entries inserted).
	Update(ctx context.Context, tmpl *domain.ProgramTemplate) error

	// GetByID retrieves a ProgramTemplate by its ID including all ProgramTemplateEntry records, returns domain.ErrNotFound if not found
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ProgramTemplate, error)

	// Archive sets archived_at on a ProgramTemplate, returns domain.ErrNotFound if not found
	Archive(ctx context.Context, id uuid.UUID) error

	// Unarchive clears archived_at on a ProgramTemplate, returns domain.ErrNotFound if not found
	Unarchive(ctx context.Context, id uuid.UUID) error

	// Delete removes a ProgramTemplate by ID (cascades to ProgramTemplateEntry records), returns domain.ErrNotFound if not found
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves ProgramTemplates with pagination.
	// limit: maximum number of records (1-100)
	// after: cursor for pagination (UUID string, base64-encoded)
	// includeArchived: when true, archived program templates are included in results
	// Returns: templates, next cursor (empty string if no more), has_more flag
	List(ctx context.Context, limit int, after string, includeArchived bool) ([]*domain.ProgramTemplate, string, bool, error)

	// ExistsByID checks if a ProgramTemplate with the given ID exists
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
}
