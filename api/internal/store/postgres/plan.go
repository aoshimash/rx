package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/aoshimash/rx/api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// planRepository implements PlanRepository with PostgreSQL
type planRepository struct {
	pool *pgxpool.Pool
}

// NewPlanRepository creates a new PostgreSQL Plan repository
func NewPlanRepository(pool *pgxpool.Pool) repository.PlanRepository {
	return &planRepository{pool: pool}
}

func (r *planRepository) GetByUserID(ctx context.Context, userID string) (*domain.Plan, error) {
	query := `
		SELECT id, name, notes, created_at, updated_at
		FROM plans
		WHERE user_id = $1
	`

	var plan domain.Plan
	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&plan.ID,
		&plan.Name,
		&plan.Notes,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get plan by user ID", "userID", userID, "error", err)
		return nil, err
	}

	sessions, err := r.getSessionsForPlan(ctx, plan.ID)
	if err != nil {
		return nil, err
	}
	plan.Sessions = sessions

	return &plan, nil
}

func (r *planRepository) getSessionsForPlan(ctx context.Context, planID uuid.UUID) ([]domain.PlanSession, error) {
	query := `
		SELECT id, plan_id, session_name, "order", date, source_program_id, source_session_id
		FROM plan_sessions
		WHERE plan_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, planID)
	if err != nil {
		slog.Error("Failed to get plan sessions", "planID", planID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var sessions []domain.PlanSession
	for rows.Next() {
		var sess domain.PlanSession
		err := rows.Scan(
			&sess.ID,
			&sess.PlanID,
			&sess.SessionName,
			&sess.Order,
			&sess.Date,
			&sess.SourceProgramID,
			&sess.SourceSessionID,
		)
		if err != nil {
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
		entriesBySession, err := r.getEntriesForPlanSessionsBatch(ctx, sessionIDs)
		if err != nil {
			return nil, err
		}
		for i := range sessions {
			sessions[i].Entries = entriesBySession[sessions[i].ID]
		}
	}

	return sessions, nil
}

func (r *planRepository) getEntriesForPlanSessionsBatch(ctx context.Context, sessionIDs []uuid.UUID) (map[uuid.UUID][]domain.PlanSessionEntry, error) {
	query := `
		SELECT id, session_id, "order", exercise_name, fields, notes
		FROM plan_session_entries
		WHERE session_id = ANY($1::uuid[])
		ORDER BY session_id, "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, uuidStrings(sessionIDs))
	if err != nil {
		slog.Error("Failed to batch get plan session entries", "error", err)
		return nil, err
	}
	defer rows.Close()

	result := make(map[uuid.UUID][]domain.PlanSessionEntry)
	for rows.Next() {
		var entry domain.PlanSessionEntry
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

func (r *planRepository) Create(ctx context.Context, userID string, plan *domain.Plan) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	id := uuid.New()
	if plan.ID != uuid.Nil {
		id = plan.ID
	}

	query := `
		INSERT INTO plans (id, user_id, name, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		id,
		userID,
		plan.Name,
		plan.Notes,
	).Scan(&plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		// Check for unique constraint violation (user already has a plan)
		if strings.Contains(err.Error(), "idx_plans_user_id") || strings.Contains(err.Error(), "duplicate key") {
			return &domain.DomainError{
				Code:    "CONFLICT",
				Message: "User already has a plan",
			}
		}
		slog.Error("Failed to create plan", "error", err)
		return err
	}

	plan.ID = id

	// Insert sessions and entries
	for i := range plan.Sessions {
		if err := r.insertSession(ctx, tx, id, &plan.Sessions[i]); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *planRepository) insertSession(ctx context.Context, tx pgx.Tx, planID uuid.UUID, session *domain.PlanSession) error {
	sessionID := uuid.New()
	if session.ID != uuid.Nil {
		sessionID = session.ID
	}

	var dateVal interface{}
	if session.Date != nil {
		dateVal = session.Date
	}

	query := `
		INSERT INTO plan_sessions (id, plan_id, session_name, "order", date, source_program_id, source_session_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.Exec(ctx, query,
		sessionID,
		planID,
		session.SessionName,
		session.Order,
		dateVal,
		session.SourceProgramID,
		session.SourceSessionID,
	)
	if err != nil {
		slog.Error("Failed to insert plan session", "error", err)
		return err
	}
	session.ID = sessionID
	session.PlanID = planID

	for j := range session.Entries {
		entryID := uuid.New()
		if session.Entries[j].ID != uuid.Nil {
			entryID = session.Entries[j].ID
		}

		fieldsJSON, err := marshalJSONBField(session.Entries[j].Fields)
		if err != nil {
			return err
		}

		entryQuery := `
			INSERT INTO plan_session_entries (id, session_id, "order", exercise_name, fields, notes)
			VALUES ($1, $2, $3, $4, $5, $6)
		`
		_, err = tx.Exec(ctx, entryQuery,
			entryID,
			sessionID,
			session.Entries[j].Order,
			session.Entries[j].ExerciseName,
			fieldsJSON,
			session.Entries[j].Notes,
		)
		if err != nil {
			slog.Error("Failed to insert plan session entry", "error", err)
			return err
		}
		session.Entries[j].ID = entryID
		session.Entries[j].SessionID = sessionID
	}

	return nil
}

func (r *planRepository) Update(ctx context.Context, userID string, plan *domain.Plan) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Update the plan row and get its ID
	var planID uuid.UUID
	err = tx.QueryRow(ctx, `
		UPDATE plans SET name = $2, notes = $3, updated_at = NOW()
		WHERE user_id = $1
		RETURNING id
	`, userID, plan.Name, plan.Notes).Scan(&planID)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to update plan", "userID", userID, "error", err)
		return err
	}
	plan.ID = planID

	// Delete existing sessions (entries cascade-deleted by DB)
	_, err = tx.Exec(ctx, `DELETE FROM plan_sessions WHERE plan_id = $1`, planID)
	if err != nil {
		slog.Error("Failed to delete plan sessions", "planID", planID, "error", err)
		return err
	}

	// Re-insert sessions and entries
	for i := range plan.Sessions {
		if err := r.insertSession(ctx, tx, planID, &plan.Sessions[i]); err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	// Re-read to get the updated_at timestamp
	updated, err := r.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	*plan = *updated
	return nil
}

func (r *planRepository) Delete(ctx context.Context, userID string) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM plans WHERE user_id = $1`, userID)
	if err != nil {
		slog.Error("Failed to delete plan", "userID", userID, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *planRepository) AddSessions(ctx context.Context, userID string, sessions []domain.PlanSession) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Get the plan ID for the user
	var planID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT id FROM plans WHERE user_id = $1`, userID).Scan(&planID)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		return err
	}

	// Determine the current max order
	var maxOrder int
	err = tx.QueryRow(ctx,
		`SELECT COALESCE(MAX("order"), -1) FROM plan_sessions WHERE plan_id = $1`,
		planID,
	).Scan(&maxOrder)
	if err != nil {
		return err
	}

	// Append sessions with sequential order values
	for i := range sessions {
		sessions[i].Order = maxOrder + 1 + i
		if err := r.insertSession(ctx, tx, planID, &sessions[i]); err != nil {
			return err
		}
	}

	// Touch updated_at
	_, err = tx.Exec(ctx, `UPDATE plans SET updated_at = NOW() WHERE id = $1`, planID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *planRepository) UpdateSession(ctx context.Context, userID string, session *domain.PlanSession) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Verify the session belongs to the user's plan
	var planID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT ps.plan_id FROM plan_sessions ps
		JOIN plans p ON ps.plan_id = p.id
		WHERE ps.id = $1 AND p.user_id = $2
	`, session.ID, userID).Scan(&planID)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		return err
	}

	var dateVal interface{}
	if session.Date != nil {
		dateVal = session.Date
	}

	// Update the session
	_, err = tx.Exec(ctx, `
		UPDATE plan_sessions
		SET session_name = $2, "order" = $3, date = $4, source_program_id = $5, source_session_id = $6
		WHERE id = $1
	`, session.ID, session.SessionName, session.Order, dateVal, session.SourceProgramID, session.SourceSessionID)
	if err != nil {
		slog.Error("Failed to update plan session", "sessionID", session.ID, "error", err)
		return err
	}

	// Delete existing entries and re-insert
	_, err = tx.Exec(ctx, `DELETE FROM plan_session_entries WHERE session_id = $1`, session.ID)
	if err != nil {
		return err
	}

	for j := range session.Entries {
		entryID := uuid.New()
		if session.Entries[j].ID != uuid.Nil {
			entryID = session.Entries[j].ID
		}

		fieldsJSON, err := marshalJSONBField(session.Entries[j].Fields)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO plan_session_entries (id, session_id, "order", exercise_name, fields, notes)
			VALUES ($1, $2, $3, $4, $5, $6)
		`,
			entryID,
			session.ID,
			session.Entries[j].Order,
			session.Entries[j].ExerciseName,
			fieldsJSON,
			session.Entries[j].Notes,
		)
		if err != nil {
			slog.Error("Failed to insert plan session entry", "error", err)
			return err
		}
		session.Entries[j].ID = entryID
		session.Entries[j].SessionID = session.ID
	}

	session.PlanID = planID

	// Touch updated_at
	_, err = tx.Exec(ctx, `UPDATE plans SET updated_at = NOW() WHERE id = $1`, planID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *planRepository) DeleteSession(ctx context.Context, userID string, sessionID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sessUUID, err := uuid.Parse(sessionID)
	if err != nil {
		return &domain.ValidationError{Field: "session_id", Message: "invalid UUID"}
	}

	// Verify the session belongs to the user's plan
	var planID uuid.UUID
	err = tx.QueryRow(ctx, `
		SELECT ps.plan_id FROM plan_sessions ps
		JOIN plans p ON ps.plan_id = p.id
		WHERE ps.id = $1 AND p.user_id = $2
	`, sessUUID, userID).Scan(&planID)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		return err
	}

	// Delete the session (entries cascade-deleted by DB)
	result, err := tx.Exec(ctx, `DELETE FROM plan_sessions WHERE id = $1`, sessUUID)
	if err != nil {
		slog.Error("Failed to delete plan session", "sessionID", sessionID, "error", err)
		return err
	}
	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	// Touch updated_at
	_, err = tx.Exec(ctx, `UPDATE plans SET updated_at = NOW() WHERE id = $1`, planID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
