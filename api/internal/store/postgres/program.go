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
		INSERT INTO programs (id, name, description, notes, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, id, program.Name, program.Description, program.Notes, program.Metadata).Scan(
		&program.CreatedAt, &program.UpdatedAt,
	)
	if err != nil {
		slog.Error("Failed to create program", "error", err)
		return err
	}

	program.ID = id

	if len(program.Entries) > 0 {
		if err = r.insertEntries(ctx, tx, id, program.Entries); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *programRepository) insertEntries(ctx context.Context, tx pgx.Tx, programID uuid.UUID, entries []domain.ProgramEntry) error {
	for i := range entries {
		entryID := uuid.New()
		if entries[i].ID != uuid.Nil {
			entryID = entries[i].ID
		}

		query := `
			INSERT INTO program_entries (
				id, program_id, "order", exercise_name,
				sets, reps, rpe, percent_1rm, notes, metadata
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

		_, err := tx.Exec(ctx, query,
			entryID,
			programID,
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
			slog.Error("Failed to insert program entry", "error", err)
			return err
		}

		entries[i].ID = entryID
		entries[i].ProgramID = programID
	}

	return nil
}

func (r *programRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Program, error) {
	query := `
		SELECT id, name, description, notes, metadata, created_at, updated_at
		FROM programs
		WHERE id = $1
	`

	var program domain.Program
	var metadataRaw []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&program.ID,
		&program.Name,
		&program.Description,
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

	entries, err := r.getEntriesForProgram(ctx, id)
	if err != nil {
		return nil, err
	}
	program.Entries = entries

	return &program, nil
}

func (r *programRepository) getEntriesForProgram(ctx context.Context, programID uuid.UUID) ([]domain.ProgramEntry, error) {
	query := `
		SELECT id, program_id, "order", exercise_name,
		       sets, reps, rpe, percent_1rm, notes, metadata
		FROM program_entries
		WHERE program_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, programID)
	if err != nil {
		slog.Error("Failed to get program entries", "programID", programID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var entries []domain.ProgramEntry
	for rows.Next() {
		var entry domain.ProgramEntry
		var metadataRaw []byte
		err := rows.Scan(
			&entry.ID,
			&entry.ProgramID,
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

func (r *programRepository) Update(ctx context.Context, program *domain.Program) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		UPDATE programs
		SET name = $2, description = $3, notes = $4, metadata = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err = tx.QueryRow(ctx, query, program.ID, program.Name, program.Description, program.Notes, program.Metadata).Scan(&program.UpdatedAt)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to update program", "id", program.ID, "error", err)
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM program_entries WHERE program_id = $1`, program.ID)
	if err != nil {
		slog.Error("Failed to delete program entries", "error", err)
		return err
	}

	if len(program.Entries) > 0 {
		if err = r.insertEntries(ctx, tx, program.ID, program.Entries); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
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

func (r *programRepository) List(ctx context.Context, limit int, after string) ([]*domain.Program, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, name, description, notes, metadata, created_at, updated_at
		FROM programs
		WHERE ($1::uuid IS NULL OR id > $1)
		ORDER BY id ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, startID, limit+1)
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
			&program.Name,
			&program.Description,
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
		entries, err := r.getEntriesForProgram(ctx, program.ID)
		if err != nil {
			return nil, "", false, err
		}
		program.Entries = entries
	}

	var nextCursor string
	if hasMore && len(programs) > 0 {
		nextCursor = encodeCursor(programs[len(programs)-1].ID)
	}

	return programs, nextCursor, hasMore, nil
}
