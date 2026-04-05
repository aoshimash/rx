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

// uuidStrings converts a []uuid.UUID to []string for use with ANY($1::uuid[]) queries.
func uuidStrings(ids []uuid.UUID) []string {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = id.String()
	}
	return strs
}

// marshalJSONBField marshals a value to JSON bytes for JSONB columns.
// Returns nil if the value is nil or empty slice/map.
func marshalJSONBField(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	if string(b) == "null" || string(b) == "[]" || string(b) == "{}" {
		return nil, nil
	}
	return b, nil
}

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
		INSERT INTO programs (id, name, notes, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		id,
		program.Name,
		program.Notes,
		string(program.Status),
	).Scan(&program.CreatedAt, &program.UpdatedAt)
	if err != nil {
		slog.Error("Failed to create program", "error", err)
		return err
	}

	program.ID = id

	// Insert groups
	for i := range program.Groups {
		groupID := uuid.New()
		if program.Groups[i].ID != uuid.Nil {
			groupID = program.Groups[i].ID
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO program_groups (id, program_id, parent_group_id, name, "order", notes)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			groupID,
			id,
			program.Groups[i].ParentGroupID,
			program.Groups[i].Name,
			program.Groups[i].Order,
			program.Groups[i].Notes,
		)
		if err != nil {
			slog.Error("Failed to insert program group", "error", err)
			return err
		}
		program.Groups[i].ID = groupID
		program.Groups[i].ProgramID = id
	}

	for i := range program.Sessions {
		sessionID := uuid.New()
		sessionQuery := `
			INSERT INTO program_sessions (id, program_id, group_id, field_group_id, session_name, "order", date)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`
		var dateVal interface{}
		if program.Sessions[i].Date != nil {
			dateVal = program.Sessions[i].Date
		}
		_, err = tx.Exec(ctx, sessionQuery,
			sessionID,
			id,
			program.Sessions[i].GroupID,
			program.Sessions[i].FieldGroupID,
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
			fieldsJSON, err := marshalJSONBField(program.Sessions[i].Entries[j].Fields)
			if err != nil {
				return err
			}
			entryQuery := `
				INSERT INTO program_session_entries (
					id, session_id, "order", exercise_name,
					fields, notes
				)
				VALUES ($1, $2, $3, $4, $5, $6)
			`
			_, err = tx.Exec(ctx, entryQuery,
				entryID,
				sessionID,
				program.Sessions[i].Entries[j].Order,
				program.Sessions[i].Entries[j].ExerciseName,
				fieldsJSON,
				program.Sessions[i].Entries[j].Notes,
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
		SELECT id, name, notes, status, created_at, updated_at
		FROM programs
		WHERE id = $1
	`

	var program domain.Program
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&program.ID,
		&program.Name,
		&program.Notes,
		&program.Status,
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

	groups, err := r.getGroupsForProgram(ctx, id)
	if err != nil {
		return nil, err
	}
	program.Groups = groups

	sessions, err := r.getSessionsForProgram(ctx, id)
	if err != nil {
		return nil, err
	}
	program.Sessions = sessions

	return &program, nil
}

func (r *programRepository) getGroupsForProgram(ctx context.Context, programID uuid.UUID) ([]domain.ProgramGroup, error) {
	query := `
		SELECT id, program_id, parent_group_id, name, "order", notes
		FROM program_groups
		WHERE program_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, programID)
	if err != nil {
		slog.Error("Failed to get program groups", "programID", programID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var groups []domain.ProgramGroup
	for rows.Next() {
		var g domain.ProgramGroup
		err := rows.Scan(
			&g.ID,
			&g.ProgramID,
			&g.ParentGroupID,
			&g.Name,
			&g.Order,
			&g.Notes,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, g)
	}

	return groups, rows.Err()
}

func (r *programRepository) getSessionsForProgram(ctx context.Context, programID uuid.UUID) ([]domain.ProgramSession, error) {
	query := `
		SELECT id, program_id, group_id, field_group_id, session_name, "order", date
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
			&sess.GroupID,
			&sess.FieldGroupID,
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
		       fields, notes
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
		var fieldsRaw []byte
		err := rows.Scan(
			&entry.ID,
			&entry.SessionID,
			&entry.Order,
			&entry.ExerciseName,
			&fieldsRaw,
			&entry.Notes,
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
		UPDATE programs SET name = $2, notes = $3, status = $4, updated_at = NOW()
		WHERE id = $1
	`, program.ID, program.Name, program.Notes, string(program.Status))
	if err != nil {
		slog.Error("Failed to update program", "id", program.ID, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	// Delete existing groups and sessions (entries cascade-deleted by DB)
	_, err = tx.Exec(ctx, `DELETE FROM program_groups WHERE program_id = $1`, program.ID)
	if err != nil {
		slog.Error("Failed to delete program groups", "id", program.ID, "error", err)
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM program_sessions WHERE program_id = $1`, program.ID)
	if err != nil {
		slog.Error("Failed to delete program sessions", "id", program.ID, "error", err)
		return err
	}

	// Re-insert groups
	for i := range program.Groups {
		groupID := uuid.New()
		if program.Groups[i].ID != uuid.Nil {
			groupID = program.Groups[i].ID
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO program_groups (id, program_id, parent_group_id, name, "order", notes)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			groupID,
			program.ID,
			program.Groups[i].ParentGroupID,
			program.Groups[i].Name,
			program.Groups[i].Order,
			program.Groups[i].Notes,
		)
		if err != nil {
			slog.Error("Failed to insert program group", "error", err)
			return err
		}
		program.Groups[i].ID = groupID
		program.Groups[i].ProgramID = program.ID
	}

	// Re-insert sessions and entries
	for i := range program.Sessions {
		sessionID := uuid.New()
		var dateVal interface{}
		if program.Sessions[i].Date != nil {
			dateVal = program.Sessions[i].Date
		}
		_, err = tx.Exec(ctx, `
			INSERT INTO program_sessions (id, program_id, group_id, field_group_id, session_name, "order", date)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, sessionID, program.ID, program.Sessions[i].GroupID, program.Sessions[i].FieldGroupID, program.Sessions[i].SessionName, program.Sessions[i].Order, dateVal)
		if err != nil {
			slog.Error("Failed to insert program session", "error", err)
			return err
		}
		program.Sessions[i].ID = sessionID
		program.Sessions[i].ProgramID = program.ID

		for j := range program.Sessions[i].Entries {
			entryID := uuid.New()
			fieldsJSON, err := marshalJSONBField(program.Sessions[i].Entries[j].Fields)
			if err != nil {
				return err
			}
			_, err = tx.Exec(ctx, `
				INSERT INTO program_session_entries (
					id, session_id, "order", exercise_name,
					fields, notes
				)
				VALUES ($1, $2, $3, $4, $5, $6)
			`,
				entryID,
				sessionID,
				program.Sessions[i].Entries[j].Order,
				program.Sessions[i].Entries[j].ExerciseName,
				fieldsJSON,
				program.Sessions[i].Entries[j].Notes,
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

func (r *programRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProgramStatus) (*domain.Program, error) {
	result, err := r.pool.Exec(ctx, `
		UPDATE programs SET status = $2, updated_at = NOW()
		WHERE id = $1
	`, id, string(status))
	if err != nil {
		slog.Error("Failed to update program status", "id", id, "error", err)
		return nil, err
	}
	if result.RowsAffected() == 0 {
		return nil, domain.ErrNotFound
	}

	return r.GetByID(ctx, id)
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
		SELECT id, name, notes, status, created_at, updated_at
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
		err := rows.Scan(
			&program.ID,
			&program.Name,
			&program.Notes,
			&program.Status,
			&program.CreatedAt,
			&program.UpdatedAt,
		)
		if err != nil {
			return nil, "", false, err
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

	if len(programs) > 0 {
		programIDs := make([]uuid.UUID, len(programs))
		for i, p := range programs {
			programIDs[i] = p.ID
		}

		groupsByProgram, err := r.getGroupsForProgramsBatch(ctx, programIDs)
		if err != nil {
			return nil, "", false, err
		}
		sessionsByProgram, err := r.getSessionsForProgramsBatch(ctx, programIDs)
		if err != nil {
			return nil, "", false, err
		}
		for _, program := range programs {
			program.Groups = groupsByProgram[program.ID]
			program.Sessions = sessionsByProgram[program.ID]
		}
	}

	var nextCursor string
	if hasMore && len(programs) > 0 {
		nextCursor = encodeCursor(programs[len(programs)-1].ID)
	}

	return programs, nextCursor, hasMore, nil
}

func (r *programRepository) getGroupsForProgramsBatch(ctx context.Context, programIDs []uuid.UUID) (map[uuid.UUID][]domain.ProgramGroup, error) {
	query := `
		SELECT id, program_id, parent_group_id, name, "order", notes
		FROM program_groups
		WHERE program_id = ANY($1::uuid[])
		ORDER BY program_id, "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, uuidStrings(programIDs))
	if err != nil {
		slog.Error("Failed to batch get program groups", "error", err)
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]domain.ProgramGroup)
	for rows.Next() {
		var g domain.ProgramGroup
		if err := rows.Scan(&g.ID, &g.ProgramID, &g.ParentGroupID, &g.Name, &g.Order, &g.Notes); err != nil {
			return nil, err
		}
		result[g.ProgramID] = append(result[g.ProgramID], g)
	}
	return result, rows.Err()
}

func (r *programRepository) getSessionsForProgramsBatch(ctx context.Context, programIDs []uuid.UUID) (map[uuid.UUID][]domain.ProgramSession, error) {
	query := `
		SELECT id, program_id, group_id, field_group_id, session_name, "order", date
		FROM program_sessions
		WHERE program_id = ANY($1::uuid[])
		ORDER BY program_id, "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, uuidStrings(programIDs))
	if err != nil {
		slog.Error("Failed to batch get program sessions", "error", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.ProgramSession
	for rows.Next() {
		var sess domain.ProgramSession
		if err := rows.Scan(&sess.ID, &sess.ProgramID, &sess.GroupID, &sess.FieldGroupID, &sess.SessionName, &sess.Order, &sess.Date); err != nil {
			return nil, err
		}
		sessions = append(sessions, sess)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(sessions) > 0 {
		sessionIDs := make([]uuid.UUID, len(sessions))
		for i, s := range sessions {
			sessionIDs[i] = s.ID
		}
		entriesBySession, err := r.getEntriesForSessionsBatch(ctx, sessionIDs)
		if err != nil {
			return nil, err
		}
		for i := range sessions {
			sessions[i].Entries = entriesBySession[sessions[i].ID]
		}
	}

	result := make(map[uuid.UUID][]domain.ProgramSession)
	for _, sess := range sessions {
		result[sess.ProgramID] = append(result[sess.ProgramID], sess)
	}
	return result, nil
}

func (r *programRepository) getEntriesForSessionsBatch(ctx context.Context, sessionIDs []uuid.UUID) (map[uuid.UUID][]domain.ProgramSessionEntry, error) {
	query := `
		SELECT id, session_id, "order", exercise_name, fields, notes
		FROM program_session_entries
		WHERE session_id = ANY($1::uuid[])
		ORDER BY session_id, "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, uuidStrings(sessionIDs))
	if err != nil {
		slog.Error("Failed to batch get program session entries", "error", err)
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]domain.ProgramSessionEntry)
	for rows.Next() {
		var entry domain.ProgramSessionEntry
		var fieldsRaw []byte
		if err := rows.Scan(&entry.ID, &entry.SessionID, &entry.Order, &entry.ExerciseName, &fieldsRaw, &entry.Notes); err != nil {
			return nil, err
		}
		if len(fieldsRaw) > 0 {
			if err := json.Unmarshal(fieldsRaw, &entry.Fields); err != nil {
				return nil, fmt.Errorf("unmarshal fields for entry %s: %w", entry.ID, err)
			}
		}
		result[entry.SessionID] = append(result[entry.SessionID], entry)
	}
	return result, rows.Err()
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
