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

	query := `
		INSERT INTO plans (id, program_id, name, description, notes, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query, id, plan.ProgramID, plan.Name, plan.Description, plan.Notes, plan.Metadata).Scan(
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

		var dateVal interface{}
		if entries[i].Date != nil {
			dateVal = time.Time(*entries[i].Date)
		}

		query := `
			INSERT INTO plan_entries (
				id, plan_id, "order", date, exercise_name,
				sets, reps, load_kg, rpe, notes, metadata
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`

		_, err := tx.Exec(ctx, query,
			entryID,
			planID,
			entries[i].Order,
			dateVal,
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
		SELECT id, program_id, name, description, notes, metadata, created_at, updated_at
		FROM plans
		WHERE id = $1
	`

	var plan domain.Plan
	var metadataRaw []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&plan.ID,
		&plan.ProgramID,
		&plan.Name,
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
		SELECT id, plan_id, "order", date, exercise_name,
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
		var dateVal *time.Time
		err := rows.Scan(
			&entry.ID,
			&entry.PlanID,
			&entry.Order,
			&dateVal,
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
		if dateVal != nil {
			d := domain.DateOnly(*dateVal)
			entry.Date = &d
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

	query := `
		UPDATE plans
		SET program_id = $2, name = $3, description = $4, notes = $5, metadata = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err = tx.QueryRow(ctx, query, plan.ID, plan.ProgramID, plan.Name, plan.Description, plan.Notes, plan.Metadata).Scan(&plan.UpdatedAt)
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
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, program_id, name, description, notes, metadata, created_at, updated_at
		FROM plans
		WHERE ($1::uuid IS NULL OR id > $1)
		ORDER BY id ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, startID, limit+1)
	if err != nil {
		slog.Error("Failed to list plans", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	plans := make([]*domain.Plan, 0, limit)

	for rows.Next() {
		var plan domain.Plan
		var metadataRaw []byte
		err := rows.Scan(
			&plan.ID,
			&plan.ProgramID,
			&plan.Name,
			&plan.Description,
			&plan.Notes,
			&metadataRaw,
			&plan.CreatedAt,
			&plan.UpdatedAt,
		)
		if err != nil {
			return nil, "", false, err
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
		nextCursor = encodeCursor(plans[len(plans)-1].ID)
	}

	return plans, nextCursor, hasMore, nil
}
