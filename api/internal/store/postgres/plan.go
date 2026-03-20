package postgres

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

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

func (r *planRepository) Create(ctx context.Context, plan *domain.Plan) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	id := uuid.New()
	if plan.ID != uuid.Nil {
		id = plan.ID
	}

	var dateVal interface{}
	if plan.Date != nil {
		dateVal = time.Time(*plan.Date)
	}

	query := `
		INSERT INTO plans (id, program_id, name, date, session_name, description, notes, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, id, plan.ProgramID, plan.Name, dateVal, plan.SessionName, plan.Description, plan.Notes, plan.Metadata).Scan(
		&plan.CreatedAt, &plan.UpdatedAt,
	)
	if err != nil {
		slog.Error("Failed to create plan", "error", err)
		return err
	}

	plan.ID = id

	if len(plan.Entries) > 0 {
		if err = r.insertEntries(ctx, tx, id, plan.Entries); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *planRepository) insertEntries(ctx context.Context, tx pgx.Tx, planID uuid.UUID, entries []domain.PlanEntry) error {
	for i := range entries {
		entryID := uuid.New()
		if entries[i].ID != uuid.Nil {
			entryID = entries[i].ID
		}

		query := `
			INSERT INTO plan_entries (
				id, plan_id, "order", exercise_name,
				sets, reps, load_kg, rpe, notes, metadata
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`

		_, err := tx.Exec(ctx, query,
			entryID,
			planID,
			entries[i].Order,
			entries[i].ExerciseName,
			entries[i].Sets,
			entries[i].Reps,
			entries[i].LoadKg,
			entries[i].RPE,
			entries[i].Notes,
			entries[i].Metadata,
		)
		if err != nil {
			slog.Error("Failed to insert plan entry", "error", err)
			return err
		}

		entries[i].ID = entryID
		entries[i].PlanID = planID
	}

	return nil
}

func (r *planRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	query := `
		SELECT id, program_id, name, date, session_name, description, notes, metadata, created_at, updated_at
		FROM plans
		WHERE id = $1
	`

	var plan domain.Plan
	var metadataRaw []byte
	var dateVal *time.Time
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&plan.ID,
		&plan.ProgramID,
		&plan.Name,
		&dateVal,
		&plan.SessionName,
		&plan.Description,
		&plan.Notes,
		&metadataRaw,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get plan by ID", "id", id, "error", err)
		return nil, err
	}

	if dateVal != nil {
		d := domain.DateOnly(*dateVal)
		plan.Date = &d
	}
	if len(metadataRaw) > 0 {
		plan.Metadata = json.RawMessage(metadataRaw)
	}

	entries, err := r.getEntriesForPlan(ctx, id)
	if err != nil {
		return nil, err
	}
	plan.Entries = entries

	return &plan, nil
}

func (r *planRepository) getEntriesForPlan(ctx context.Context, planID uuid.UUID) ([]domain.PlanEntry, error) {
	query := `
		SELECT id, plan_id, "order", exercise_name,
		       sets, reps, load_kg, rpe, notes, metadata
		FROM plan_entries
		WHERE plan_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, planID)
	if err != nil {
		slog.Error("Failed to get plan entries", "planID", planID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var entries []domain.PlanEntry
	for rows.Next() {
		var entry domain.PlanEntry
		var metadataRaw []byte
		err := rows.Scan(
			&entry.ID,
			&entry.PlanID,
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

func (r *planRepository) Update(ctx context.Context, plan *domain.Plan) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	var updateDateVal interface{}
	if plan.Date != nil {
		updateDateVal = time.Time(*plan.Date)
	}

	query := `
		UPDATE plans
		SET program_id = $2, name = $3, date = $4, session_name = $5,
		    description = $6, notes = $7, metadata = $8, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err = tx.QueryRow(ctx, query, plan.ID, plan.ProgramID, plan.Name, updateDateVal, plan.SessionName, plan.Description, plan.Notes, plan.Metadata).Scan(&plan.UpdatedAt)
	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to update plan", "id", plan.ID, "error", err)
		return err
	}

	_, err = tx.Exec(ctx, `DELETE FROM plan_entries WHERE plan_id = $1`, plan.ID)
	if err != nil {
		slog.Error("Failed to delete plan entries", "error", err)
		return err
	}

	if len(plan.Entries) > 0 {
		if err = r.insertEntries(ctx, tx, plan.ID, plan.Entries); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *planRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM plans WHERE id = $1`, id)
	if err != nil {
		slog.Error("Failed to delete plan", "id", id, "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *planRepository) List(ctx context.Context, limit int, after string) ([]*domain.Plan, string, bool, error) {
	var startTime *time.Time
	var startID uuid.UUID
	if after != "" {
		t, id, err := decodePlanCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		startTime = &t
		startID = id
	}

	query := `
		SELECT id, program_id, name, date, session_name, description, notes, metadata, created_at, updated_at
		FROM plans
		WHERE NOT EXISTS (SELECT 1 FROM logs WHERE logs.plan_id = plans.id)
		  AND ($1::timestamptz IS NULL OR (created_at, id) > ($1, $2::uuid))
		ORDER BY created_at ASC, id ASC
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, startTime, startID, limit+1)
	if err != nil {
		slog.Error("Failed to list plans", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	return r.scanPlanList(ctx, rows, limit)
}

func (r *planRepository) ListByProgramID(ctx context.Context, programID uuid.UUID, limit int, after string) ([]*domain.Plan, string, bool, error) {
	var startTime *time.Time
	var startID uuid.UUID
	if after != "" {
		t, id, err := decodePlanCursor(after)
		if err != nil {
			return nil, "", false, err
		}
		startTime = &t
		startID = id
	}

	query := `
		SELECT id, program_id, name, date, session_name, description, notes, metadata, created_at, updated_at
		FROM plans
		WHERE program_id = $4
		  AND NOT EXISTS (SELECT 1 FROM logs WHERE logs.plan_id = plans.id)
		  AND ($1::timestamptz IS NULL OR (created_at, id) > ($1, $2::uuid))
		ORDER BY created_at ASC, id ASC
		LIMIT $3
	`

	rows, err := r.pool.Query(ctx, query, startTime, startID, limit+1, programID)
	if err != nil {
		slog.Error("Failed to list plans", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	return r.scanPlanList(ctx, rows, limit)
}

func (r *planRepository) scanPlanList(ctx context.Context, rows pgx.Rows, limit int) ([]*domain.Plan, string, bool, error) {
	plans := make([]*domain.Plan, 0, limit)

	for rows.Next() {
		var plan domain.Plan
		var metadataRaw []byte
		var dateVal *time.Time
		err := rows.Scan(
			&plan.ID,
			&plan.ProgramID,
			&plan.Name,
			&dateVal,
			&plan.SessionName,
			&plan.Description,
			&plan.Notes,
			&metadataRaw,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}
		if dateVal != nil {
			d := domain.DateOnly(*dateVal)
			plan.Date = &d
		}
		if len(metadataRaw) > 0 {
			plan.Metadata = json.RawMessage(metadataRaw)
		}
		plans = append(plans, &plan)
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(plans) > limit
	if hasMore {
		plans = plans[:limit]
	}

	for _, plan := range plans {
		entries, err := r.getEntriesForPlan(ctx, plan.ID)
		if err != nil {
			return nil, "", false, err
		}
		plan.Entries = entries
	}

	var nextCursor string
	if hasMore && len(plans) > 0 {
		last := plans[len(plans)-1]
		nextCursor = encodePlanCursor(last.CreatedAt, last.ID)
	}

	return plans, nextCursor, hasMore, nil
}

// encodePlanCursor encodes a (created_at, id) pair for plan cursor-based pagination
func encodePlanCursor(createdAt time.Time, id uuid.UUID) string {
	s := createdAt.Format(time.RFC3339Nano) + "|" + id.String()
	return base64.URLEncoding.EncodeToString([]byte(s))
}

// decodePlanCursor decodes a plan cursor to (created_at, id)
func decodePlanCursor(cursor string) (time.Time, uuid.UUID, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	parts := strings.SplitN(string(data), "|", 2)
	if len(parts) != 2 {
		return time.Time{}, uuid.Nil, fmt.Errorf("invalid plan cursor format")
	}
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	id, err := uuid.Parse(parts[1])
	if err != nil {
		return time.Time{}, uuid.Nil, err
	}
	return t, id, nil
}
