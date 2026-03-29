package postgres

import (
	"context"
	"encoding/json"
	"fmt"
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
		INSERT INTO logs (id, program_id, session_name, performed_at, started_at, finished_at, notes, metadata, plan_snapshot, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
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
		log.PlanSnapshot,
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

		fieldsJSON, err := marshalJSONBField(entries[i].Fields)
		if err != nil {
			return err
		}

		query := `
			INSERT INTO log_entries (
				id, log_id, "order", exercise_name,
				fields,
				notes, video_object_key, started_at, finished_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`

		_, err = tx.Exec(ctx, query,
			entryID,
			logID,
			entries[i].Order,
			entries[i].ExerciseName,
			fieldsJSON,
			entries[i].Notes,
			entries[i].VideoObjectKey,
			entries[i].StartedAt,
			entries[i].FinishedAt,
		)
		if err != nil {
			slog.Error("Failed to insert log entry", "error", err)
			return err
		}

		entries[i].ID = entryID
		entries[i].LogID = logID

		// Insert sets for this entry
		for k := range entries[i].Sets {
			setID := uuid.New()
			if entries[i].Sets[k].ID != uuid.Nil {
				setID = entries[i].Sets[k].ID
			}

			setFieldsJSON, err := marshalJSONBField(entries[i].Sets[k].Fields)
			if err != nil {
				return err
			}

			setQuery := `
				INSERT INTO log_sets (id, entry_id, set_number, fields, video_url, notes)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
			_, err = tx.Exec(ctx, setQuery,
				setID,
				entryID,
				entries[i].Sets[k].SetNumber,
				setFieldsJSON,
				entries[i].Sets[k].VideoURL,
				entries[i].Sets[k].Notes,
			)
			if err != nil {
				slog.Error("Failed to insert log set", "error", err)
				return err
			}
			entries[i].Sets[k].ID = setID
			entries[i].Sets[k].EntryID = entryID
		}
	}

	return nil
}

func (r *logRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Log, error) {
	query := `
		SELECT id, program_id, session_name, performed_at, started_at, finished_at, notes, metadata, plan_snapshot, created_at, updated_at
		FROM logs
		WHERE id = $1
	`

	var log domain.Log
	var metadataRaw []byte
	var planSnapshotRaw []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&log.ID,
		&log.ProgramID,
		&log.SessionName,
		&log.PerformedAt,
		&log.StartedAt,
		&log.FinishedAt,
		&log.Notes,
		&metadataRaw,
		&planSnapshotRaw,
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
	if len(planSnapshotRaw) > 0 {
		log.PlanSnapshot = json.RawMessage(planSnapshotRaw)
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
		       fields,
		       notes, video_object_key, started_at, finished_at
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
		var fieldsRaw []byte
		err := rows.Scan(
			&entry.ID,
			&entry.LogID,
			&entry.Order,
			&entry.ExerciseName,
			&fieldsRaw,
			&entry.Notes,
			&entry.VideoObjectKey,
			&entry.StartedAt,
			&entry.FinishedAt,
		)
		if err != nil {
			return nil, err
		}
		if len(fieldsRaw) > 0 {
			if err := json.Unmarshal(fieldsRaw, &entry.Fields); err != nil {
				return nil, fmt.Errorf("unmarshal fields for entry %s: %w", entry.ID, err)
			}
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load sets for each entry
	for i := range entries {
		sets, err := r.getSetsForEntry(ctx, entries[i].ID)
		if err != nil {
			return nil, err
		}
		entries[i].Sets = sets
	}

	return entries, nil
}

func (r *logRepository) getSetsForEntry(ctx context.Context, entryID uuid.UUID) ([]domain.LogSet, error) {
	query := `
		SELECT id, entry_id, set_number, fields, video_url, notes
		FROM log_sets
		WHERE entry_id = $1
		ORDER BY set_number ASC
	`

	rows, err := r.pool.Query(ctx, query, entryID)
	if err != nil {
		slog.Error("Failed to get log sets", "entryID", entryID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var sets []domain.LogSet
	for rows.Next() {
		var s domain.LogSet
		var fieldsRaw []byte
		err := rows.Scan(
			&s.ID,
			&s.EntryID,
			&s.SetNumber,
			&fieldsRaw,
			&s.VideoURL,
			&s.Notes,
		)
		if err != nil {
			return nil, err
		}
		if len(fieldsRaw) > 0 {
			if err := json.Unmarshal(fieldsRaw, &s.Fields); err != nil {
				return nil, fmt.Errorf("unmarshal fields for set %s: %w", s.ID, err)
			}
		}
		sets = append(sets, s)
	}

	return sets, rows.Err()
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
		    notes = $7, metadata = $8, plan_snapshot = $9, updated_at = NOW()
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
		log.PlanSnapshot,
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
		SELECT id, program_id, session_name, performed_at, started_at, finished_at, notes, metadata, plan_snapshot, created_at, updated_at
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
		var planSnapshotRaw []byte
		err := rows.Scan(
			&log.ID,
			&log.ProgramID,
			&log.SessionName,
			&log.PerformedAt,
			&log.StartedAt,
			&log.FinishedAt,
			&log.Notes,
			&metadataRaw,
			&planSnapshotRaw,
			&log.CreatedAt,
			&log.UpdatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}
		if len(metadataRaw) > 0 {
			log.Metadata = json.RawMessage(metadataRaw)
		}
		if len(planSnapshotRaw) > 0 {
			log.PlanSnapshot = json.RawMessage(planSnapshotRaw)
		}
		logs = append(logs, &log)
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(logs) > limit
	if hasMore {
		logs = logs[:limit]
	}

	if len(logs) > 0 {
		logIDs := make([]uuid.UUID, len(logs))
		for i, l := range logs {
			logIDs[i] = l.ID
		}
		entriesByLog, err := r.getEntriesForLogsBatch(ctx, logIDs)
		if err != nil {
			return nil, "", false, err
		}
		for _, l := range logs {
			l.Entries = entriesByLog[l.ID]
		}
	}

	var nextCursor string
	if hasMore && len(logs) > 0 {
		nextCursor = encodeCursor(logs[len(logs)-1].ID)
	}

	return logs, nextCursor, hasMore, nil
}

func (r *logRepository) getEntriesForLogsBatch(ctx context.Context, logIDs []uuid.UUID) (map[uuid.UUID][]domain.LogEntry, error) {
	query := `
		SELECT id, log_id, "order", exercise_name,
		       fields,
		       notes, video_object_key, started_at, finished_at
		FROM log_entries
		WHERE log_id = ANY($1::uuid[])
		ORDER BY log_id, "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, uuidStrings(logIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []domain.LogEntry
	for rows.Next() {
		var entry domain.LogEntry
		var fieldsRaw []byte
		if err := rows.Scan(
			&entry.ID, &entry.LogID, &entry.Order, &entry.ExerciseName,
			&fieldsRaw, &entry.Notes, &entry.VideoObjectKey, &entry.StartedAt, &entry.FinishedAt,
		); err != nil {
			return nil, err
		}
		if len(fieldsRaw) > 0 {
			if err := json.Unmarshal(fieldsRaw, &entry.Fields); err != nil {
				return nil, fmt.Errorf("unmarshal fields for entry %s: %w", entry.ID, err)
			}
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(entries) > 0 {
		entryIDs := make([]uuid.UUID, len(entries))
		for i, e := range entries {
			entryIDs[i] = e.ID
		}
		setsByEntry, err := r.getSetsForEntriesBatch(ctx, entryIDs)
		if err != nil {
			return nil, err
		}
		for i := range entries {
			entries[i].Sets = setsByEntry[entries[i].ID]
		}
	}

	result := make(map[uuid.UUID][]domain.LogEntry)
	for _, entry := range entries {
		result[entry.LogID] = append(result[entry.LogID], entry)
	}
	return result, nil
}

func (r *logRepository) getSetsForEntriesBatch(ctx context.Context, entryIDs []uuid.UUID) (map[uuid.UUID][]domain.LogSet, error) {
	query := `
		SELECT id, entry_id, set_number, fields, video_url, notes
		FROM log_sets
		WHERE entry_id = ANY($1::uuid[])
		ORDER BY entry_id, set_number ASC
	`

	rows, err := r.pool.Query(ctx, query, uuidStrings(entryIDs))
	if err != nil {
		slog.Error("Failed to batch get log sets", "error", err)
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]domain.LogSet)
	for rows.Next() {
		var s domain.LogSet
		var fieldsRaw []byte
		if err := rows.Scan(&s.ID, &s.EntryID, &s.SetNumber, &fieldsRaw, &s.VideoURL, &s.Notes); err != nil {
			return nil, err
		}
		if len(fieldsRaw) > 0 {
			if err := json.Unmarshal(fieldsRaw, &s.Fields); err != nil {
				return nil, fmt.Errorf("unmarshal fields for set %s: %w", s.ID, err)
			}
		}
		result[s.EntryID] = append(result[s.EntryID], s)
	}
	return result, rows.Err()
}
