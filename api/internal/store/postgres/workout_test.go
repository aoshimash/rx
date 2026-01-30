package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/google/uuid"
)

// stringPtr is defined in exercise_test.go

func TestWorkoutRepository_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	// Create exercise first (required for workout entry)
	exerciseRepo := NewExerciseRepository(pool)
	exercise := &domain.Exercise{Name: "Test Exercise"}
	if err := exerciseRepo.Create(ctx, exercise); err != nil {
		t.Fatalf("Failed to create exercise: %v", err)
	}

	repo := NewWorkoutRepository(pool)

	workout := &domain.Workout{
		Timestamp: time.Now(),
		Entries: []domain.WorkoutEntry{
			{
				ExerciseID: exercise.ID,
				EntryType:  stringPtr("work"),
				Sets:       3,
				Reps:       10,
				LoadKg:     100.0,
				RPE:        8,
				Order:      0,
			},
		},
	}

	err = repo.Create(ctx, workout)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if workout.ID == uuid.Nil {
		t.Error("Create() did not set ID")
	}
	if len(workout.Entries) != 1 {
		t.Errorf("Create() Entries length = %v, want 1", len(workout.Entries))
	}
	if workout.Entries[0].ID == uuid.Nil {
		t.Error("Create() did not set Entry ID")
	}
}

func TestWorkoutRepository_GetByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	exerciseRepo := NewExerciseRepository(pool)
	exercise := &domain.Exercise{Name: "Test Exercise"}
	if err := exerciseRepo.Create(ctx, exercise); err != nil {
		t.Fatalf("Failed to create exercise: %v", err)
	}

	repo := NewWorkoutRepository(pool)

	workout := &domain.Workout{
		Timestamp: time.Now(),
		Entries: []domain.WorkoutEntry{
			{
				ExerciseID: exercise.ID,
				EntryType:  stringPtr("work"),
				Sets:       3,
				Reps:       10,
				LoadKg:     100.0,
				RPE:        8,
				Order:      0,
			},
		},
	}
	if err := repo.Create(ctx, workout); err != nil {
		t.Fatalf("Failed to create workout: %v", err)
	}

	got, err := repo.GetByID(ctx, workout.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.ID != workout.ID {
		t.Errorf("GetByID() ID = %v, want %v", got.ID, workout.ID)
	}
	if len(got.Entries) != 1 {
		t.Errorf("GetByID() Entries length = %v, want 1", len(got.Entries))
	}
}

func TestWorkoutRepository_Update(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	exerciseRepo := NewExerciseRepository(pool)
	exercise := &domain.Exercise{Name: "Test Exercise"}
	if err := exerciseRepo.Create(ctx, exercise); err != nil {
		t.Fatalf("Failed to create exercise: %v", err)
	}

	repo := NewWorkoutRepository(pool)

	workout := &domain.Workout{
		Timestamp: time.Now(),
		Entries: []domain.WorkoutEntry{
			{
				ExerciseID: exercise.ID,
				EntryType:  stringPtr("work"),
				Sets:       3,
				Reps:       10,
				LoadKg:     100.0,
				RPE:        8,
				Order:      0,
			},
		},
	}
	if err := repo.Create(ctx, workout); err != nil {
		t.Fatalf("Failed to create workout: %v", err)
	}

	workout.Entries[0].Sets = 5
	if err := repo.Update(ctx, workout); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	got, err := repo.GetByID(ctx, workout.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.Entries[0].Sets != 5 {
		t.Errorf("GetByID() Sets = %v, want 5", got.Entries[0].Sets)
	}
}

func TestWorkoutRepository_Delete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	exerciseRepo := NewExerciseRepository(pool)
	exercise := &domain.Exercise{Name: "Test Exercise"}
	if err := exerciseRepo.Create(ctx, exercise); err != nil {
		t.Fatalf("Failed to create exercise: %v", err)
	}

	repo := NewWorkoutRepository(pool)

	workout := &domain.Workout{
		Timestamp: time.Now(),
		Entries: []domain.WorkoutEntry{
			{
				ExerciseID: exercise.ID,
				EntryType:  stringPtr("work"),
				Sets:       3,
				Reps:       10,
				LoadKg:     100.0,
				RPE:        8,
				Order:      0,
			},
		},
	}
	if err := repo.Create(ctx, workout); err != nil {
		t.Fatalf("Failed to create workout: %v", err)
	}

	if err := repo.Delete(ctx, workout.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(ctx, workout.ID)
	if err != domain.ErrNotFound {
		t.Errorf("GetByID() after Delete() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestWorkoutRepository_List(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	exerciseRepo := NewExerciseRepository(pool)
	exercise := &domain.Exercise{Name: "Test Exercise"}
	if err := exerciseRepo.Create(ctx, exercise); err != nil {
		t.Fatalf("Failed to create exercise: %v", err)
	}

	repo := NewWorkoutRepository(pool)

	// Create multiple workouts
	for i := 0; i < 3; i++ {
		workout := &domain.Workout{
			Timestamp: time.Now().Add(time.Duration(i) * time.Hour),
			Entries: []domain.WorkoutEntry{
				{
					ExerciseID: exercise.ID,
					EntryType:  stringPtr("work"),
					Sets:       3,
					Reps:       10,
					LoadKg:     100.0,
					RPE:        8,
					Order:      0,
				},
			},
		}
		if err := repo.Create(ctx, workout); err != nil {
			t.Fatalf("Failed to create workout: %v", err)
		}
	}

	workouts, _, hasMore, err := repo.List(ctx, 2, "")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(workouts) != 2 {
		t.Errorf("List() returned %d workouts, want 2", len(workouts))
	}
	if !hasMore {
		t.Error("List() hasMore = false, want true")
	}
}
