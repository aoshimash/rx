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

// cycleRepository implements CycleRepository with PostgreSQL
type cycleRepository struct {
	pool *pgxpool.Pool
}

// NewCycleRepository creates a new PostgreSQL Cycle repository
func NewCycleRepository(pool *pgxpool.Pool) repository.CycleRepository {
	return &cycleRepository{pool: pool}
}

func (r *cycleRepository) Create(ctx context.Context, cycle *domain.Cycle) error {
	id := uuid.New()
	if cycle.ID != uuid.Nil {
		id = cycle.ID
	}

	query := `
		INSERT INTO cycles (id, program_id, name, notes, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING created_at
	`

	err := r.pool.QueryRow(ctx, query, id, cycle.ProgramID, cycle.Name, cycle.Notes, cycle.Metadata).Scan(
		&cycle.CreatedAt,
	)
	if err != nil {
		slog.Error("Failed to create cycle", "error", err)
		return err
	}

	cycle.ID = id
	return nil
}

func (r *cycleRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Cycle, error) {
	query := `
		SELECT id, program_id, name, notes, metadata, created_at
		FROM cycles
		WHERE id = $1
	`

	var cycle domain.Cycle
	var metadataRaw []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&cycle.ID,
		&cycle.ProgramID,
		&cycle.Name,
		&cycle.Notes,
		&metadataRaw,
		&cycle.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get cycle by ID", "id", id, "error", err)
		return nil, err
	}

	if len(metadataRaw) > 0 {
		cycle.Metadata = json.RawMessage(metadataRaw)
	}

	return &cycle, nil
}

func (r *cycleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM cycles WHERE id = $1`, id)
	if err != nil {
		slog.Error("Failed to delete cycle", "id", id, "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *cycleRepository) List(ctx context.Context, limit int, after string) ([]*domain.Cycle, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, program_id, name, notes, metadata, created_at
		FROM cycles
		WHERE ($1::uuid IS NULL OR id > $1)
		ORDER BY id ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, startID, limit+1)
	if err != nil {
		slog.Error("Failed to list cycles", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	return r.scanCycleList(rows, limit)
}

func (r *cycleRepository) ListByProgramID(ctx context.Context, programID uuid.UUID, limit int, after string) ([]*domain.Cycle, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, program_id, name, notes, metadata, created_at
		FROM cycles
		WHERE program_id = $3
		  AND ($1::uuid IS NULL OR id > $1)
		ORDER BY id ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, startID, limit+1, programID)
	if err != nil {
		slog.Error("Failed to list cycles by program ID", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	return r.scanCycleList(rows, limit)
}

func (r *cycleRepository) ExistsByProgramID(ctx context.Context, programID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM cycles WHERE program_id = $1)`, programID).Scan(&exists)
	if err != nil {
		slog.Error("Failed to check cycles by program ID", "error", err)
		return false, err
	}
	return exists, nil
}

func (r *cycleRepository) scanCycleList(rows pgx.Rows, limit int) ([]*domain.Cycle, string, bool, error) {
	cycles := make([]*domain.Cycle, 0, limit)

	for rows.Next() {
		var cycle domain.Cycle
		var metadataRaw []byte
		err := rows.Scan(
			&cycle.ID,
			&cycle.ProgramID,
			&cycle.Name,
			&cycle.Notes,
			&metadataRaw,
			&cycle.CreatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}
		if len(metadataRaw) > 0 {
			cycle.Metadata = json.RawMessage(metadataRaw)
		}
		cycles = append(cycles, &cycle)
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(cycles) > limit
	if hasMore {
		cycles = cycles[:limit]
	}

	var nextCursor string
	if hasMore && len(cycles) > 0 {
		nextCursor = encodeCursor(cycles[len(cycles)-1].ID)
	}

	return cycles, nextCursor, hasMore, nil
}
