package postgres

import (
	"context"
	"encoding/base64"
	"log/slog"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/aoshimash/optel-training/api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// exerciseRepository implements ExerciseRepository with PostgreSQL
type exerciseRepository struct {
	pool *pgxpool.Pool
}

// NewExerciseRepository creates a new PostgreSQL Exercise repository
func NewExerciseRepository(pool *pgxpool.Pool) repository.ExerciseRepository {
	return &exerciseRepository{pool: pool}
}

func (r *exerciseRepository) Create(ctx context.Context, exercise *domain.Exercise) error {
	query := `
		INSERT INTO exercises (id, name, description, aliases, muscle_groups, load_increment, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	id := uuid.New()
	if exercise.ID != uuid.Nil {
		id = exercise.ID
	}

	err := r.pool.QueryRow(ctx, query,
		id,
		exercise.Name,
		exercise.Description,
		exercise.Aliases,
		exercise.MuscleGroups,
		exercise.LoadIncrement,
	).Scan(&exercise.CreatedAt, &exercise.UpdatedAt)

	if err != nil {
		slog.Error("Failed to create exercise", "error", err)
		return err
	}

	exercise.ID = id
	return nil
}

func (r *exerciseRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Exercise, error) {
	query := `
		SELECT id, name, description, aliases, muscle_groups, load_increment, created_at, updated_at
		FROM exercises
		WHERE id = $1
	`

	var exercise domain.Exercise
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.Aliases,
		&exercise.MuscleGroups,
		&exercise.LoadIncrement,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get exercise by ID", "id", id, "error", err)
		return nil, err
	}

	return &exercise, nil
}

func (r *exerciseRepository) Update(ctx context.Context, exercise *domain.Exercise) error {
	query := `
		UPDATE exercises
		SET name = $2, description = $3, aliases = $4, muscle_groups = $5, load_increment = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err := r.pool.QueryRow(ctx, query,
		exercise.ID,
		exercise.Name,
		exercise.Description,
		exercise.Aliases,
		exercise.MuscleGroups,
		exercise.LoadIncrement,
	).Scan(&exercise.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to update exercise", "id", exercise.ID, "error", err)
		return err
	}

	return nil
}

func (r *exerciseRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM exercises WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete exercise", "id", id, "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *exerciseRepository) List(ctx context.Context, limit int, after string) ([]*domain.Exercise, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, name, description, aliases, muscle_groups, load_increment, created_at, updated_at
		FROM exercises
		WHERE ($1::uuid IS NULL OR id > $1)
		ORDER BY id ASC
		LIMIT $2
	`

	rows, err := r.pool.Query(ctx, query, startID, limit+1)
	if err != nil {
		slog.Error("Failed to list exercises", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	exercises := make([]*domain.Exercise, 0, limit)
	var lastID uuid.UUID

	for rows.Next() {
		var exercise domain.Exercise
		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.Aliases,
			&exercise.MuscleGroups,
			&exercise.LoadIncrement,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}

		exercises = append(exercises, &exercise)
		lastID = exercise.ID
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(exercises) > limit
	if hasMore {
		exercises = exercises[:limit]
	}

	var nextCursor string
	if hasMore && len(exercises) > 0 {
		nextCursor = encodeCursor(lastID)
	}

	return exercises, nextCursor, hasMore, nil
}

// encodeCursor encodes a UUID as a base64 string for pagination
func encodeCursor(id uuid.UUID) string {
	return base64.URLEncoding.EncodeToString(id[:])
}

// decodeCursor decodes a base64 string to a UUID
func decodeCursor(cursor string) (uuid.UUID, error) {
	data, err := base64.URLEncoding.DecodeString(cursor)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.FromBytes(data)
}
