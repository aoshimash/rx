package postgres

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// logRepository implements LogRepository with PostgreSQL
type logRepository struct {
	pool *pgxpool.Pool
}

// NewLogRepository creates a new PostgreSQL Log repository
func NewLogRepository(pool *pgxpool.Pool) repository.LogRepository {
	return &logRepository{pool: pool}
}

func (r *logRepository) Create(ctx context.Context, log *domain.Log) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	id := uuid.New()
	if log.ID != uuid.Nil {
		id = log.ID
	}

	query := `
		INSERT INTO logs (id, program_id, session_name, performed_at, started_at, finished_at, notes, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		id,
		log.ProgramID,
		log.SessionName,
		log.PerformedAt,
		log.StartedAt,
		log.FinishedAt,
		log.Notes,
		log.Metadata,
	).Scan(&log.CreatedAt, &log.UpdatedAt)

	if err != nil {
		slog.Error("Failed to create log", "error", err)
		return err
	}

	log.ID = id

	if len(log.Entries) > 0 {
		err = r.insertEntries(ctx, tx, id, log.Entries)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *logRepository) insertEntries(ctx context.Context, tx pgx.Tx, logID uuid.UUID, entries []domain.LogEntry) error {
	for i := range entries {
		entryID := uuid.New()
		if entries[i].ID != uuid.Nil {
			entryID = entries[i].ID
		}

		query := `
			INSERT INTO log_entries (
				id, log_id, "order", exercise_name,
				sets, reps, load_kg, rpe,
				notes, video_object_key, started_at, finished_at, metadata
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`

		_, err := tx.Exec(ctx, query,
			entryID,
			logID,
			entries[i].Order,
			entries[i].ExerciseName,
			entries[i].Sets,
			entries[i].Reps,
			entries[i].LoadKg,
			entries[i].RPE,
			entries[i].Notes,
			entries[i].VideoObjectKey,
			entries[i].StartedAt,
			entries[i].FinishedAt,
			entries[i].Metadata,
		)
		if err != nil {
			slog.Error("Failed to insert log entry", "error", err)
			return err
		}

		entries[i].ID = entryID
		entries[i].LogID = logID
	}

	return nil
}

func (r *logRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Log, error) {
	query := `
		SELECT id, program_id, session_name, performed_at, started_at, finished_at, notes, metadata, created_at, updated_at
		FROM logs
		WHERE id = $1
	`

	var log domain.Log
	var metadataRaw []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&log.ID,
		&log.ProgramID,
		&log.SessionName,
		&log.PerformedAt,
		&log.StartedAt,
		&log.FinishedAt,
		&log.Notes,
		&metadataRaw,
		&log.CreatedAt,
		&log.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get log by ID", "id", id, "error", err)
		return nil, err
	}

	if len(metadataRaw) > 0 {
		log.Metadata = json.RawMessage(metadataRaw)
	}

	entries, err := r.getEntriesForLog(ctx, id)
	if err != nil {
		slog.Error("Failed to get log entries", "id", id, "error", err)
		return nil, err
	}

	log.Entries = entries
	return &log, nil
}

func (r *logRepository) getEntriesForLog(ctx context.Context, logID uuid.UUID) ([]domain.LogEntry, error) {
	query := `
		SELECT id, log_id, "order", exercise_name,
		       sets, reps, load_kg, rpe,
		       notes, video_object_key, started_at, finished_at, metadata
		FROM log_entries
		WHERE log_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]domain.LogEntry, 0)
	for rows.Next() {
		var entry domain.LogEntry
		var metadataRaw []byte
		err := rows.Scan(
			&entry.ID,
			&entry.LogID,
			&entry.Order,
			&entry.ExerciseName,
			&entry.Sets,
			&entry.Reps,
			&entry.LoadKg,
			&entry.RPE,
			&entry.Notes,
			&entry.VideoObjectKey,
			&entry.StartedAt,
			&entry.FinishedAt,
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

func (r *logRepository) Update(ctx context.Context, log *domain.Log) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	query := `
		UPDATE logs
		SET program_id = $2, session_name = $3, performed_at = $4, started_at = $5, finished_at = $6,
		    notes = $7, metadata = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err = tx.QueryRow(ctx, query,
		log.ID,
		log.ProgramID,
		log.SessionName,
		log.PerformedAt,
		log.StartedAt,
		log.FinishedAt,
		log.Notes,
		log.Metadata,
	).Scan(&log.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to update log", "id", log.ID, "error", err)
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM log_entries WHERE log_id = $1`, log.ID)
	if err != nil {
		slog.Error("Failed to delete log entries", "error", err)
		return err
	}

	if len(log.Entries) > 0 {
		err = r.insertEntries(ctx, tx, log.ID, log.Entries)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *logRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM logs WHERE id = $1`, id)
	if err != nil {
		slog.Error("Failed to delete log", "id", id, "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *logRepository) List(ctx context.Context, programID *uuid.UUID, limit int, after string) ([]*domain.Log, string, bool, error) {
	return r.listWithFilter(ctx, programID, nil, nil, limit, after)
}

func (r *logRepository) ListByPerformedAtRange(ctx context.Context, programID *uuid.UUID, performedAtFrom, performedAtTo *time.Time, limit int, after string) ([]*domain.Log, string, bool, error) {
	return r.listWithFilter(ctx, programID, performedAtFrom, performedAtTo, limit, after)
}

func (r *logRepository) listWithFilter(ctx context.Context, programID *uuid.UUID, performedAtFrom, performedAtTo *time.Time, limit int, after string) ([]*domain.Log, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, program_id, session_name, performed_at, started_at, finished_at, notes, metadata, created_at, updated_at
		FROM logs
		WHERE ($1::uuid IS NULL OR id > $1)
		  AND ($2::uuid IS NULL OR program_id = $2)
		  AND ($3::timestamptz IS NULL OR performed_at >= $3)
		  AND ($4::timestamptz IS NULL OR performed_at < $4)
		ORDER BY id ASC
		LIMIT $5
	`

	rows, err := r.pool.Query(ctx, query, startID, programID, performedAtFrom, performedAtTo, limit+1)
	if err != nil {
		slog.Error("Failed to list logs", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	logs := make([]*domain.Log, 0, limit)

	for rows.Next() {
		var log domain.Log
		var metadataRaw []byte
		err := rows.Scan(
			&log.ID,
			&log.ProgramID,
			&log.SessionName,
			&log.PerformedAt,
			&log.StartedAt,
			&log.FinishedAt,
			&log.Notes,
			&metadataRaw,
			&log.CreatedAt,
			&log.UpdatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}
		if len(metadataRaw) > 0 {
			log.Metadata = json.RawMessage(metadataRaw)
		}

		entries, err := r.getEntriesForLog(ctx, log.ID)
		if err != nil {
			return nil, "", false, err
		}
		log.Entries = entries

		logs = append(logs, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(logs) > limit
	if hasMore {
		logs = logs[:limit]
	}

	var nextCursor string
	if hasMore && len(logs) > 0 {
		nextCursor = encodeCursor(logs[len(logs)-1].ID)
	}

	return logs, nextCursor, hasMore, nil
}

func (r *logRepository) ListLoggedSessionsByProgramID(ctx context.Context, programID uuid.UUID) ([]domain.LoggedSession, error) {
	query := `
		SELECT session_name, id
		FROM logs
		WHERE program_id = $1 AND session_name IS NOT NULL
		ORDER BY session_name
	`

	rows, err := r.pool.Query(ctx, query, programID)
	if err != nil {
		slog.Error("Failed to list logged sessions", "programID", programID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.LoggedSession
	for rows.Next() {
		var s domain.LoggedSession
		if err := rows.Scan(&s.SessionName, &s.LogID); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	if sessions == nil {
		sessions = []domain.LoggedSession{}
	}

	return sessions, rows.Err()
}

func (r *logRepository) ExistsByProgramIDAndSessionName(ctx context.Context, programID uuid.UUID, sessionName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM logs WHERE program_id = $1 AND session_name = $2)`
	var exists bool
	if err := r.pool.QueryRow(ctx, query, programID, sessionName).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *logRepository) ExistsByProgramIDAndSessionNameExcluding(ctx context.Context, programID uuid.UUID, sessionName string, excludeID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM logs WHERE program_id = $1 AND session_name = $2 AND id != $3)`
	var exists bool
	if err := r.pool.QueryRow(ctx, query, programID, sessionName, excludeID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
