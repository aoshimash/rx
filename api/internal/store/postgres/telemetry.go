package postgres

import (
	"context"
	"log/slog"
	"time"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/aoshimash/optel-workout/api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// telemetryRepository implements TelemetryPointRepository with PostgreSQL
type telemetryRepository struct {
	pool *pgxpool.Pool
}

// NewTelemetryPointRepository creates a new PostgreSQL TelemetryPoint repository
func NewTelemetryPointRepository(pool *pgxpool.Pool) repository.TelemetryPointRepository {
	return &telemetryRepository{pool: pool}
}

func (r *telemetryRepository) Create(ctx context.Context, point *domain.TelemetryPoint) error {
	query := `
		INSERT INTO telemetry_points (id, timestamp, metric_name, value, unit, workout_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING created_at
	`

	id := uuid.New()
	if point.ID != uuid.Nil {
		id = point.ID
	}

	err := r.pool.QueryRow(ctx, query,
		id,
		point.Timestamp,
		point.MetricName,
		point.Value,
		point.Unit,
		point.WorkoutID,
	).Scan(&point.CreatedAt)

	if err != nil {
		slog.Error("Failed to create telemetry point", "error", err)
		return err
	}

	point.ID = id
	return nil
}

func (r *telemetryRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.TelemetryPoint, error) {
	query := `
		SELECT id, timestamp, metric_name, value, unit, workout_id, created_at
		FROM telemetry_points
		WHERE id = $1
	`

	var point domain.TelemetryPoint
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&point.ID,
		&point.Timestamp,
		&point.MetricName,
		&point.Value,
		&point.Unit,
		&point.WorkoutID,
		&point.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get telemetry point by ID", "id", id, "error", err)
		return nil, err
	}

	return &point, nil
}

func (r *telemetryRepository) Update(ctx context.Context, point *domain.TelemetryPoint) error {
	query := `
		UPDATE telemetry_points
		SET timestamp = $2, metric_name = $3, value = $4, unit = $5, workout_id = $6
		WHERE id = $1
	`

	result, err := r.pool.Exec(ctx, query,
		point.ID,
		point.Timestamp,
		point.MetricName,
		point.Value,
		point.Unit,
		point.WorkoutID,
	)

	if err != nil {
		slog.Error("Failed to update telemetry point", "id", point.ID, "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *telemetryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM telemetry_points WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete telemetry point", "id", id, "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *telemetryRepository) List(ctx context.Context, limit int, after string) ([]*domain.TelemetryPoint, string, bool, error) {
	return r.listWithFilter(ctx, "", nil, nil, limit, after)
}

func (r *telemetryRepository) ListByMetricAndTimeRange(ctx context.Context, metricName string, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.TelemetryPoint, string, bool, error) {
	return r.listWithFilter(ctx, metricName, timestampFrom, timestampTo, limit, after)
}

func (r *telemetryRepository) listWithFilter(ctx context.Context, metricName string, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.TelemetryPoint, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, timestamp, metric_name, value, unit, workout_id, created_at
		FROM telemetry_points
		WHERE ($1::uuid IS NULL OR id > $1)
		  AND ($2::text IS NULL OR metric_name = $2)
		  AND ($3::timestamptz IS NULL OR timestamp >= $3)
		  AND ($4::timestamptz IS NULL OR timestamp < $4)
		ORDER BY id ASC
		LIMIT $5
	`

	var metricNameParam *string
	if metricName != "" {
		metricNameParam = &metricName
	}

	rows, err := r.pool.Query(ctx, query, startID, metricNameParam, timestampFrom, timestampTo, limit+1)
	if err != nil {
		slog.Error("Failed to list telemetry points", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	points := make([]*domain.TelemetryPoint, 0, limit)

	for rows.Next() {
		var point domain.TelemetryPoint
		err := rows.Scan(
			&point.ID,
			&point.Timestamp,
			&point.MetricName,
			&point.Value,
			&point.Unit,
			&point.WorkoutID,
			&point.CreatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}

		points = append(points, &point)
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(points) > limit
	if hasMore {
		points = points[:limit]
	}

	var nextCursor string
	if hasMore && len(points) > 0 {
		// Use the last item in the returned set, not the extra item
		nextCursor = encodeCursor(points[len(points)-1].ID)
	}

	return points, nextCursor, hasMore, nil
}
