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

// workoutRepository implements WorkoutRepository with PostgreSQL
type workoutRepository struct {
	pool *pgxpool.Pool
}

// NewWorkoutRepository creates a new PostgreSQL Workout repository
func NewWorkoutRepository(pool *pgxpool.Pool) repository.WorkoutRepository {
	return &workoutRepository{pool: pool}
}

// scanWorkoutRow scans a workout row from the database
func scanWorkoutRow(row pgx.Row) (*domain.Workout, error) {
	var workout domain.Workout
	err := row.Scan(
		&workout.ID,
		&workout.Timestamp,
		&workout.SessionStart,
		&workout.SessionEnd,
		&workout.BodyWeightKg,
		&workout.FatigueLevel,
		&workout.SleepHours,
		&workout.ConditionNotes,
		&workout.ProgramNodeID,
		&workout.ProgramContext,
		&workout.Notes,
		&workout.CreatedAt,
		&workout.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

// scanWorkoutEntryRow scans a workout entry row from the database
func scanWorkoutEntryRow(row pgx.Row) (*domain.WorkoutEntry, []byte, error) {
	var entry domain.WorkoutEntry
	var planSnapshotJSON []byte

	err := row.Scan(
		&entry.ID,
		&entry.WorkoutID,
		&entry.Order,
		&entry.ExerciseID,
		&entry.DisplayName,
		&entry.EntryType,
		&entry.Sets,
		&entry.Reps,
		&entry.LoadKg,
		&entry.RPE,
		&entry.EntryStart,
		&entry.EntryEnd,
		&entry.PlannedRestSeconds,
		&entry.PerformedRestSeconds,
		&entry.PerSetRestOverrides,
		&entry.ProgramNodeID,
		&planSnapshotJSON,
		&entry.Notes,
		&entry.VideoObjectKey,
	)
	if err != nil {
		return nil, nil, err
	}

	// Deserialize PlanSnapshot from JSONB
	if len(planSnapshotJSON) > 0 {
		var planSnapshot domain.PlanSnapshot
		if err := json.Unmarshal(planSnapshotJSON, &planSnapshot); err != nil {
			return nil, nil, err
		}
		entry.PlanSnapshot = &planSnapshot
	}

	return &entry, planSnapshotJSON, nil
}

func (r *workoutRepository) Create(ctx context.Context, workout *domain.Workout) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	id := uuid.New()
	if workout.ID != uuid.Nil {
		id = workout.ID
	}

	// Insert workout
	query := `
		INSERT INTO workouts (
			id, timestamp, session_start, session_end,
			body_weight_kg, fatigue_level, sleep_hours,
			condition_notes, program_node_id, program_context,
			notes, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING created_at, updated_at
	`

	err = tx.QueryRow(ctx, query,
		id,
		workout.Timestamp,
		workout.SessionStart,
		workout.SessionEnd,
		workout.BodyWeightKg,
		workout.FatigueLevel,
		workout.SleepHours,
		workout.ConditionNotes,
		workout.ProgramNodeID,
		workout.ProgramContext,
		workout.Notes,
	).Scan(&workout.CreatedAt, &workout.UpdatedAt)

	if err != nil {
		slog.Error("Failed to create workout", "error", err)
		return err
	}

	workout.ID = id

	// Insert workout entries
	if len(workout.Entries) > 0 {
		err = r.insertEntries(ctx, tx, id, workout.Entries)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *workoutRepository) insertEntries(ctx context.Context, tx pgx.Tx, workoutID uuid.UUID, entries []domain.WorkoutEntry) error {
	for i := range entries {
		entryID := uuid.New()
		if entries[i].ID != uuid.Nil {
			entryID = entries[i].ID
		}

		// Serialize PlanSnapshot to JSONB
		var planSnapshotJSON []byte
		if entries[i].PlanSnapshot != nil {
			var err error
			planSnapshotJSON, err = json.Marshal(entries[i].PlanSnapshot)
			if err != nil {
				return err
			}
		}

		query := `
			INSERT INTO workout_entries (
				id, workout_id, "order", exercise_id, display_name,
				entry_type, sets, reps, load_kg, rpe,
				entry_start, entry_end, planned_rest_seconds,
				performed_rest_seconds, per_set_rest_overrides,
				program_node_id, plan_snapshot, notes, video_object_key
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		`

		_, err := tx.Exec(ctx, query,
			entryID,
			workoutID,
			entries[i].Order,
			entries[i].ExerciseID,
			entries[i].DisplayName,
			entries[i].EntryType,
			entries[i].Sets,
			entries[i].Reps,
			entries[i].LoadKg,
			entries[i].RPE,
			entries[i].EntryStart,
			entries[i].EntryEnd,
			entries[i].PlannedRestSeconds,
			entries[i].PerformedRestSeconds,
			entries[i].PerSetRestOverrides,
			entries[i].ProgramNodeID,
			planSnapshotJSON,
			entries[i].Notes,
			entries[i].VideoObjectKey,
		)
		if err != nil {
			slog.Error("Failed to insert workout entry", "error", err)
			return err
		}

		entries[i].ID = entryID
		entries[i].WorkoutID = workoutID
	}

	return nil
}

func (r *workoutRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Workout, error) {
	// Get workout
	query := `
		SELECT id, timestamp, session_start, session_end,
		       body_weight_kg, fatigue_level, sleep_hours,
		       condition_notes, program_node_id, program_context,
		       notes, created_at, updated_at
		FROM workouts
		WHERE id = $1
	`

	workout, err := scanWorkoutRow(r.pool.QueryRow(ctx, query, id))
	if err == pgx.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to get workout by ID", "id", id, "error", err)
		return nil, err
	}

	// Get workout entries
	entries, err := r.getEntriesForWorkout(ctx, id)
	if err != nil {
		slog.Error("Failed to get workout entries", "id", id, "error", err)
		return nil, err
	}

	workout.Entries = entries
	return workout, nil
}

func (r *workoutRepository) Update(ctx context.Context, workout *domain.Workout) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Update workout
	query := `
		UPDATE workouts
		SET timestamp = $2, session_start = $3, session_end = $4,
		    body_weight_kg = $5, fatigue_level = $6, sleep_hours = $7,
		    condition_notes = $8, program_node_id = $9, program_context = $10,
		    notes = $11, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at
	`

	err = tx.QueryRow(ctx, query,
		workout.ID,
		workout.Timestamp,
		workout.SessionStart,
		workout.SessionEnd,
		workout.BodyWeightKg,
		workout.FatigueLevel,
		workout.SleepHours,
		workout.ConditionNotes,
		workout.ProgramNodeID,
		workout.ProgramContext,
		workout.Notes,
	).Scan(&workout.UpdatedAt)

	if err == pgx.ErrNoRows {
		return domain.ErrNotFound
	}
	if err != nil {
		slog.Error("Failed to update workout", "id", workout.ID, "error", err)
		return err
	}

	// Delete existing entries (CASCADE will handle this, but explicit for clarity)
	_, err = tx.Exec(ctx, `DELETE FROM workout_entries WHERE workout_id = $1`, workout.ID)
	if err != nil {
		slog.Error("Failed to delete workout entries", "error", err)
		return err
	}

	// Insert new entries
	if len(workout.Entries) > 0 {
		err = r.insertEntries(ctx, tx, workout.ID, workout.Entries)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *workoutRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM workouts WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		slog.Error("Failed to delete workout", "id", id, "error", err)
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *workoutRepository) List(ctx context.Context, limit int, after string) ([]*domain.Workout, string, bool, error) {
	return r.listWithFilter(ctx, nil, nil, limit, after)
}

func (r *workoutRepository) ListByTimestampRange(ctx context.Context, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.Workout, string, bool, error) {
	return r.listWithFilter(ctx, timestampFrom, timestampTo, limit, after)
}

func (r *workoutRepository) listWithFilter(ctx context.Context, timestampFrom, timestampTo *time.Time, limit int, after string) ([]*domain.Workout, string, bool, error) {
	var startID uuid.UUID
	if after != "" {
		var err error
		startID, err = decodeCursor(after)
		if err != nil {
			return nil, "", false, err
		}
	}

	query := `
		SELECT id, timestamp, session_start, session_end,
		       body_weight_kg, fatigue_level, sleep_hours,
		       condition_notes, program_node_id, program_context,
		       notes, created_at, updated_at
		FROM workouts
		WHERE ($1::uuid IS NULL OR id > $1)
		  AND ($2::timestamptz IS NULL OR timestamp >= $2)
		  AND ($3::timestamptz IS NULL OR timestamp < $3)
		ORDER BY id ASC
		LIMIT $4
	`

	rows, err := r.pool.Query(ctx, query, startID, timestampFrom, timestampTo, limit+1)
	if err != nil {
		slog.Error("Failed to list workouts", "error", err)
		return nil, "", false, err
	}
	defer rows.Close()

	workouts := make([]*domain.Workout, 0, limit)

	for rows.Next() {
		workout, err := scanWorkoutRow(rows)
		if err != nil {
			return nil, "", false, err
		}

		// Load entries for each workout
		entries, err := r.getEntriesForWorkout(ctx, workout.ID)
		if err != nil {
			return nil, "", false, err
		}
		workout.Entries = entries

		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		return nil, "", false, err
	}

	hasMore := len(workouts) > limit
	if hasMore {
		workouts = workouts[:limit]
	}

	var nextCursor string
	if hasMore && len(workouts) > 0 {
		// Use the last item in the returned set, not the extra item
		nextCursor = encodeCursor(workouts[len(workouts)-1].ID)
	}

	return workouts, nextCursor, hasMore, nil
}

func (r *workoutRepository) getEntriesForWorkout(ctx context.Context, workoutID uuid.UUID) ([]domain.WorkoutEntry, error) {
	query := `
		SELECT id, workout_id, "order", exercise_id, display_name,
		       entry_type, sets, reps, load_kg, rpe,
		       entry_start, entry_end, planned_rest_seconds,
		       performed_rest_seconds, per_set_rest_overrides,
		       program_node_id, plan_snapshot, notes, video_object_key
		FROM workout_entries
		WHERE workout_id = $1
		ORDER BY "order" ASC
	`

	rows, err := r.pool.Query(ctx, query, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]domain.WorkoutEntry, 0)
	for rows.Next() {
		entry, _, err := scanWorkoutEntryRow(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, *entry)
	}

	return entries, rows.Err()
}

func (r *workoutRepository) ListByExerciseID(ctx context.Context, exerciseID uuid.UUID) ([]*domain.Workout, error) {
	query := `
		SELECT DISTINCT w.id, w.timestamp, w.session_start, w.session_end,
		       w.body_weight_kg, w.fatigue_level, w.sleep_hours,
		       w.condition_notes, w.program_node_id, w.program_context,
		       w.notes, w.created_at, w.updated_at
		FROM workouts w
		INNER JOIN workout_entries we ON w.id = we.workout_id
		WHERE we.exercise_id = $1
		ORDER BY w.timestamp DESC
	`

	rows, err := r.pool.Query(ctx, query, exerciseID)
	if err != nil {
		slog.Error("Failed to list workouts by exercise ID", "exercise_id", exerciseID, "error", err)
		return nil, err
	}
	defer rows.Close()

	workouts := make([]*domain.Workout, 0)
	for rows.Next() {
		workout, err := scanWorkoutRow(rows)
		if err != nil {
			return nil, err
		}

		// Load entries for each workout
		entries, err := r.getEntriesForWorkout(ctx, workout.ID)
		if err != nil {
			return nil, err
		}
		workout.Entries = entries

		workouts = append(workouts, workout)
	}

	return workouts, rows.Err()
}

func (r *workoutRepository) ListByProgramNodeID(ctx context.Context, programNodeID uuid.UUID) ([]*domain.Workout, error) {
	query := `
		SELECT id, timestamp, session_start, session_end,
		       body_weight_kg, fatigue_level, sleep_hours,
		       condition_notes, program_node_id, program_context,
		       notes, created_at, updated_at
		FROM workouts
		WHERE program_node_id = $1
		ORDER BY timestamp DESC
	`

	rows, err := r.pool.Query(ctx, query, programNodeID)
	if err != nil {
		slog.Error("Failed to list workouts by program node ID", "program_node_id", programNodeID, "error", err)
		return nil, err
	}
	defer rows.Close()

	workouts := make([]*domain.Workout, 0)
	for rows.Next() {
		workout, err := scanWorkoutRow(rows)
		if err != nil {
			return nil, err
		}

		// Load entries for each workout
		entries, err := r.getEntriesForWorkout(ctx, workout.ID)
		if err != nil {
			return nil, err
		}
		workout.Entries = entries

		workouts = append(workouts, workout)
	}

	return workouts, rows.Err()
}
