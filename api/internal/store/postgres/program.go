package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// programRepository implements ProgramRepository with PostgreSQL
type programRepository struct {
	pool *pgxpool.Pool
}

// NewProgramRepository creates a new PostgreSQL Program repository
func NewProgramRepository(pool *pgxpool.Pool) repository.ProgramRepository {
	return &programRepository{pool: pool}
}

func (r *programRepository) Create(ctx context.Context, program *domain.Program) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	id := uuid.New()
	if program.ID != uuid.Nil {
		id = program.ID
	}

	query := `
		INSERT INTO programs (id, program_template_id, name, status, notes, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		id,
		program.ProgramTemplateID,
		program.Name,
		program.Status,
		program.Notes,
		program.Metadata,
	).Scan(&program.CreatedAt, &program.UpdatedAt)
	if err != nil {
		slog.Error("Failed to create program", "error", err)
		return err
	}

	program.ID = id

	for i := range program.Sessions {
		sessionID := uuid.New()
		sessionQuery := `
			INSERT INTO program_sessions (id, program_id, session_name, "order", date)
			VALUES ($1, $2, $3, $4, $5)
		`
		var dateVal interface{}
		if program.Sessions[i].Date != nil {
			dateVal = program.Sessions[i].Date
		}
		_, err = tx.Exec(ctx, sessionQuery,
			sessionID,
			id,
			program.Sessions[i].SessionName,
			program.Sessions[i].Order,
			dateVal,
		)
		if err != nil {
			slog.Error("Failed to insert program session", "error", err)
			return err
		}
		program.Sessions[i].ID = sessionID
		program.Sessions[i].ProgramID = id

		for j := range program.Sessions[i].Entries {
			entryID := uuid.New()
			entryQuery := `
				INSERT INTO program_session_entries (
					id, session_id, "order", exercise_name,
					sets, reps, load_kg, rpe, notes, metadata
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			`
			_, err = tx.Exec(ctx, entryQuery,
				entryID,
				sessionID,
				program.Sessions[i].Entries[j].Order,
				program.Sessions[i].Entries[j].ExerciseName,
				program.Sessions[i].Entries[j].Sets,
				program.Sessions[i].Entries[j].Reps,
				program.Sessions[i].Entries[j].LoadKg,
				program.Sessions[i].Entries[j].RPE,
				program.Sessions[i].Entries[j].Notes,
				program.Sessions[i].Entries[j].Metadata,
			)
			if err != nil {
				slog.Error("Failed to insert program session entry", "error", err)
				return err
			}
			program.Sessions[i].Entries[j].ID = entryID
			program.Sessions[i].Entries[j].SessionID = sessionID
		}
	}

	return tx.Commit(ctx)
}

func (r *programRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error) {
	query := `
		SELECT id, program_template_id, name, status, notes, metadata, created_at, updated_at
		FROM programs
		WHERE id = $1
	`

	var program domain.Program
	var metadataRaw []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&program.ID,
		&program.ProgramTemplateID,
		&program.Name,
		&program.Status,
		&program.Notes,
		&metadataRaw,
		&program.CreatedAt,
		&program.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get program by ID", "id", id, "error", err)
		return nil, err
	}

	if len(metadataRaw) > 0 {
		program.Metadata = json.RawMessage(metadataRaw)
	}

	sessions, err := r.getSessionsForProgram(ctx, id)
	if err != nil {
		return nil, err
	}
	program.Sessions = sessions

	return &program, nil
}

func (r *programRepository) getSessionsForProgram(ctx context.Context, programID uuid.UUID) ([]domain.ProgramSession, error) {
	query := `
		SELECT id, program_id, session_name, "order", date
		FROM program_sessions
		WHERE program_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, programID)
	if err != nil {
		slog.Error("Failed to get program sessions", "programID", programID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.ProgramSession
	for rows.Next() {
		var sess domain.ProgramSession
		err := rows.Scan(
			&sess.ID,
			&sess.ProgramID,
			&sess.SessionName,
			&sess.Order,
			&sess.Date,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, sess)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range sessions {
		entries, err := r.getEntriesForSession(ctx, sessions[i].ID)
		if err != nil {
			return nil, err
		}
		sessions[i].Entries = entries
	}

	return sessions, nil
}

func (r *programRepository) getEntriesForSession(ctx context.Context, sessionID uuid.UUID) ([]domain.ProgramSessionEntry, error) {
	query := `
		SELECT id, session_id, "order", exercise_name,
		       sets, reps, load_kg, rpe, notes, metadata
		FROM program_session_entries
		WHERE session_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, sessionID)
	if err != nil {
		slog.Error("Failed to get program session entries", "sessionID", sessionID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var entries []domain.ProgramSessionEntry
	for rows.Next() {
		var entry domain.ProgramSessionEntry
		var metadataRaw []byte
		err := rows.Scan(
			&entry.ID,
			&entry.SessionID,
			&entry.Order,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.LoadKg,
			&entry.RPE,
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

func (r *programRepository) Update(ctx context.Context, program *domain.Program) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Update the program row
	result, err := tx.Exec(ctx, `
		UPDATE programs SET name = $2, notes = $3, metadata = $4, updated_at = NOW()
		WHERE id = $1
	`, program.ID, program.Name, program.Notes, program.Metadata)
	if err != nil {
		slog.Error("Failed to update program", "id", program.ID, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	// Delete existing sessions (entries cascade-deleted by DB)
	_, err = tx.Exec(ctx, `DELETE FROM program_sessions WHERE program_id = $1`, program.ID)
	if err != nil {
		slog.Error("Failed to delete program sessions", "id", program.ID, "error", err)
		return err
	}

	// Re-insert sessions and entries
	for i := range program.Sessions {
		sessionID := uuid.New()
		var dateVal interface{}
		if program.Sessions[i].Date != nil {
			dateVal = program.Sessions[i].Date
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO program_sessions (id, program_id, session_name, "order", date)
			VALUES ($1, $2, $3, $4, $5)
		`, sessionID, program.ID, program.Sessions[i].SessionName, program.Sessions[i].Order, dateVal)
		if err != nil {
			slog.Error("Failed to insert program session", "error", err)
			return err
		}
		program.Sessions[i].ID = sessionID
		program.Sessions[i].ProgramID = program.ID

		for j := range program.Sessions[i].Entries {
			entryID := uuid.New()
			_, err = tx.Exec(ctx, `
				INSERT INTO program_session_entries (
					id, session_id, "order", exercise_name,
					sets, reps, load_kg, rpe, notes, metadata
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			`,
				entryID,
				sessionID,
				program.Sessions[i].Entries[j].Order,
				program.Sessions[i].Entries[j].ExerciseName,
				program.Sessions[i].Entries[j].Sets,
				program.Sessions[i].Entries[j].Reps,
				program.Sessions[i].Entries[j].LoadKg,
				program.Sessions[i].Entries[j].RPE,
				program.Sessions[i].Entries[j].Notes,
				program.Sessions[i].Entries[j].Metadata,
			)
			if err != nil {
				slog.Error("Failed to insert program session entry", "error", err)
				return err
			}
			program.Sessions[i].Entries[j].ID = entryID
			program.Sessions[i].Entries[j].SessionID = sessionID
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// Re-read to get the updated_at timestamp
	updated, err := r.GetByID(ctx, program.ID)
	if err != nil {
		return err
	}
	*program = *updated
	return nil
}

func (r *programRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProgramStatus) error {
	result, err := r.pool.Exec(ctx,
		`UPDATE programs SET status = $2, updated_at = NOW() WHERE id = $1`,
		id, status,
	)
	if err != nil {
		slog.Error("Failed to update program status", "id", id, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *programRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM programs WHERE id = $1`, id)
	if err != nil {
		slog.Error("Failed to delete program", "id", id, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *programRepository) List(ctx context.Context, limit int, after string, programTemplateID *uuid.UUID, status string) ([]*domain.Program, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, program_template_id, name, status, notes, metadata, created_at, updated_at
		FROM programs
		WHERE ($1::uuid IS NULL OR id > $1)
		  AND ($2::uuid IS NULL OR program_template_id = $2)
		  AND ($3::text = '' OR status = $3)
		ORDER BY id ASC
		LIMIT $4
	`

	rows, err := r.pool.Query(ctx, query, startID, programTemplateID, status, limit+1)
	if err != nil {
		slog.Error("Failed to list programs", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	programs := make([]*domain.Program, 0, limit)
	for rows.Next() {
		var program domain.Program
		var metadataRaw []byte
		err := rows.Scan(
			&program.ID,
			&program.ProgramTemplateID,
			&program.Name,
			&program.Status,
			&program.Notes,
			&metadataRaw,
			&program.CreatedAt,
			&program.UpdatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}
		if len(metadataRaw) > 0 {
			program.Metadata = json.RawMessage(metadataRaw)
		}
		programs = append(programs, &program)
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(programs) > limit
	if hasMore {
		programs = programs[:limit]
	}

	for _, program := range programs {
		sessions, err := r.getSessionsForProgram(ctx, program.ID)
		if err != nil {
			return nil, "", false, err
		}
		program.Sessions = sessions
	}

	var nextCursor string
	if hasMore && len(programs) > 0 {
		nextCursor = encodeCursor(programs[len(programs)-1].ID)
	}

	return programs, nextCursor, hasMore, nil
}

func (r *programRepository) ExistsByProgramTemplateID(ctx context.Context, programTemplateID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM programs WHERE program_template_id = $1)`,
		programTemplateID,
	).Scan(&exists)
	if err != nil {
		slog.Error("Failed to check program existence by template ID", "programTemplateID", programTemplateID, "error", err)
		return false, err
	}
	return exists, nil
}

func (r *programRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM programs WHERE name = $1)`,
		name,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check program name exists: %w", err)
	}
	return exists, nil
}
