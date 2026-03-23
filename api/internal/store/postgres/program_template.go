package postgres

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// programTemplateRepository implements ProgramTemplateRepository with PostgreSQL
type programTemplateRepository struct {
	pool *pgxpool.Pool
}

// NewProgramTemplateRepository creates a new PostgreSQL ProgramTemplate repository
func NewProgramTemplateRepository(pool *pgxpool.Pool) repository.ProgramTemplateRepository {
	return &programTemplateRepository{pool: pool}
}

func (r *programTemplateRepository) Create(ctx context.Context, tmpl *domain.ProgramTemplate) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	id := uuid.New()
	if tmpl.ID != uuid.Nil {
		id = tmpl.ID
	}

	query := `
		INSERT INTO program_templates (id, name, description, notes, metadata, weeks, days_per_week, created_by, source_template_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, id, tmpl.Name, tmpl.Description, tmpl.Notes, tmpl.Metadata, tmpl.Weeks, tmpl.DaysPerWeek, tmpl.CreatedBy, tmpl.SourceTemplateID).Scan(
		&tmpl.CreatedAt, &tmpl.UpdatedAt,
	)
	if err != nil {
		slog.Error("Failed to create program template", "error", err)
		return err
	}

	tmpl.ID = id

	if len(tmpl.Entries) > 0 {
		if err = r.insertEntries(ctx, tx, id, tmpl.Entries); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *programTemplateRepository) insertEntries(ctx context.Context, tx pgx.Tx, templateID uuid.UUID, entries []domain.ProgramTemplateEntry) error {
	for i := range entries {
		entryID := uuid.New()
		if entries[i].ID != uuid.Nil {
			entryID = entries[i].ID
		}

		query := `
			INSERT INTO program_template_entries (
				id, program_template_id, "order", exercise_name,
				sets, reps, rpe, percent_1rm, notes, metadata
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

		_, err := tx.Exec(ctx, query,
			entryID,
			templateID,
			entries[i].Order,
			entries[i].ExerciseName,
			entries[i].Sets,
			entries[i].Reps,
			entries[i].RPE,
			entries[i].Percent1RM,
			entries[i].Notes,
			entries[i].Metadata,
		)
		if err != nil {
			slog.Error("Failed to insert program template entry", "error", err)
			return err
		}

		entries[i].ID = entryID
		entries[i].ProgramTemplateID = templateID
	}

	return nil
}

func (r *programTemplateRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProgramTemplate, error) {
	query := `
		SELECT id, name, description, notes, metadata, weeks, days_per_week, created_by, created_at, updated_at, archived_at, source_template_id
		FROM program_templates
		WHERE id = $1
	`

	var tmpl domain.ProgramTemplate
	var metadataRaw []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&tmpl.ID,
		&tmpl.Name,
		&tmpl.Description,
		&tmpl.Notes,
		&metadataRaw,
		&tmpl.Weeks,
		&tmpl.DaysPerWeek,
		&tmpl.CreatedBy,
		&tmpl.CreatedAt,
		&tmpl.UpdatedAt,
		&tmpl.ArchivedAt,
		&tmpl.SourceTemplateID,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get program template by ID", "id", id, "error", err)
		return nil, err
	}

	if len(metadataRaw) > 0 {
		tmpl.Metadata = json.RawMessage(metadataRaw)
	}

	entries, err := r.getEntriesForTemplate(ctx, id)
	if err != nil {
		return nil, err
	}
	tmpl.Entries = entries

	return &tmpl, nil
}

func (r *programTemplateRepository) getEntriesForTemplate(ctx context.Context, templateID uuid.UUID) ([]domain.ProgramTemplateEntry, error) {
	query := `
		SELECT id, program_template_id, "order", exercise_name,
		       sets, reps, rpe, percent_1rm, notes, metadata
		FROM program_template_entries
		WHERE program_template_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, templateID)
	if err != nil {
		slog.Error("Failed to get program template entries", "templateID", templateID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var entries []domain.ProgramTemplateEntry
	for rows.Next() {
		var entry domain.ProgramTemplateEntry
		var metadataRaw []byte
		err := rows.Scan(
			&entry.ID,
			&entry.ProgramTemplateID,
			&entry.Order,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.RPE,
			&entry.Percent1RM,
			&entry.Notes,
			&metadataRaw,
		)
		if err != nil {
			return nil, err
		}
		if len(metadataRaw) > 0 {
			entry.Metadata = json.RawMessage(metadataRaw)
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func (r *programTemplateRepository) CreateAndArchive(ctx context.Context, tmpl *domain.ProgramTemplate, archiveID uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Create new template
	id := uuid.New()
	if tmpl.ID != uuid.Nil {
		id = tmpl.ID
	}

	query := `
		INSERT INTO program_templates (id, name, description, notes, metadata, weeks, days_per_week, created_by, source_template_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, id, tmpl.Name, tmpl.Description, tmpl.Notes, tmpl.Metadata, tmpl.Weeks, tmpl.DaysPerWeek, tmpl.CreatedBy, tmpl.SourceTemplateID).Scan(
		&tmpl.CreatedAt, &tmpl.UpdatedAt,
	)
	if err != nil {
		slog.Error("Failed to create program template in CreateAndArchive", "error", err)
		return err
	}

	tmpl.ID = id

	if len(tmpl.Entries) > 0 {
		if err = r.insertEntries(ctx, tx, id, tmpl.Entries); err != nil {
			return err
		}
	}

	// Archive old template
	result, err := tx.Exec(ctx,
		`UPDATE program_templates SET archived_at = NOW() WHERE id = $1`,
		archiveID,
	)
	if err != nil {
		slog.Error("Failed to archive program template in CreateAndArchive", "id", archiveID, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return tx.Commit(ctx)
}

func (r *programTemplateRepository) Archive(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE program_templates SET archived_at = NOW() WHERE id = $1`,
		id,
	)
	if err != nil {
		slog.Error("Failed to archive program template", "id", id, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *programTemplateRepository) Unarchive(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE program_templates SET archived_at = NULL WHERE id = $1`,
		id,
	)
	if err != nil {
		slog.Error("Failed to unarchive program template", "id", id, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *programTemplateRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM program_templates WHERE id = $1`, id)
	if err != nil {
		slog.Error("Failed to delete program template", "id", id, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *programTemplateRepository) List(ctx context.Context, limit int, after string, includeArchived bool) ([]*domain.ProgramTemplate, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	var query string
	if includeArchived {
		query = `
			SELECT id, name, description, notes, metadata, weeks, days_per_week, created_by, created_at, updated_at, archived_at, source_template_id
			FROM program_templates
			WHERE ($1::uuid IS NULL OR id > $1)
			ORDER BY id ASC
			LIMIT $2
		`
	} else {
		query = `
			SELECT id, name, description, notes, metadata, weeks, days_per_week, created_by, created_at, updated_at, archived_at, source_template_id
			FROM program_templates
			WHERE ($1::uuid IS NULL OR id > $1) AND archived_at IS NULL
			ORDER BY id ASC
			LIMIT $2
		`
	}

	rows, err := r.pool.Query(ctx, query, startID, limit+1)
	if err != nil {
		slog.Error("Failed to list program templates", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	templates := make([]*domain.ProgramTemplate, 0, limit)
	for rows.Next() {
		var tmpl domain.ProgramTemplate
		var metadataRaw []byte
		err := rows.Scan(
			&tmpl.ID,
			&tmpl.Name,
			&tmpl.Description,
			&tmpl.Notes,
			&metadataRaw,
			&tmpl.Weeks,
			&tmpl.DaysPerWeek,
			&tmpl.CreatedBy,
			&tmpl.CreatedAt,
			&tmpl.UpdatedAt,
			&tmpl.ArchivedAt,
			&tmpl.SourceTemplateID,
		)
		if err != nil {
			return nil, "", false, err
		}
		if len(metadataRaw) > 0 {
			tmpl.Metadata = json.RawMessage(metadataRaw)
		}
		templates = append(templates, &tmpl)
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(templates) > limit
	if hasMore {
		templates = templates[:limit]
	}

	for _, tmpl := range templates {
		entries, err := r.getEntriesForTemplate(ctx, tmpl.ID)
		if err != nil {
			return nil, "", false, err
		}
		tmpl.Entries = entries
	}

	var nextCursor string
	if hasMore && len(templates) > 0 {
		nextCursor = encodeCursor(templates[len(templates)-1].ID)
	}

	return templates, nextCursor, hasMore, nil
}

func (r *programTemplateRepository) Update(ctx context.Context, tmpl *domain.ProgramTemplate) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	_, err = tx.Exec(ctx, `DELETE FROM program_template_entries WHERE program_template_id = $1`, tmpl.ID)
	if err != nil {
		slog.Error("Failed to delete old program template entries", "id", tmpl.ID, "error", err)
		return err
	}

	query := `
		UPDATE program_templates
		SET name = $2, description = $3, notes = $4, metadata = $5,
		    weeks = $6, days_per_week = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at, source_template_id
	`
	err = tx.QueryRow(ctx, query,
		tmpl.ID, tmpl.Name, tmpl.Description, tmpl.Notes, tmpl.Metadata,
		tmpl.Weeks, tmpl.DaysPerWeek,
	).Scan(&tmpl.UpdatedAt, &tmpl.SourceTemplateID)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to update program template", "id", tmpl.ID, "error", err)
		return err
	}

	if len(tmpl.Entries) > 0 {
		if err = r.insertEntries(ctx, tx, tmpl.ID, tmpl.Entries); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *programTemplateRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM program_templates WHERE id = $1)`,
		id,
	).Scan(&exists)
	if err != nil {
		slog.Error("Failed to check program template existence", "id", id, "error", err)
		return false, err
	}
	return exists, nil
}
