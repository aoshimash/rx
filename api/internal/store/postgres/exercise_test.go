package postgres

import (
	"context"
	"testing"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/google/uuid"
)

func TestExerciseRepository_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewExerciseRepository(pool)

	tests := []struct {
		name     string
		exercise *domain.Exercise
		wantErr  bool
	}{
		{
			name: "create valid exercise",
			exercise: &domain.Exercise{
				Name: "Bench Press",
			},
			wantErr: false,
		},
		{
			name: "create exercise with all fields",
			exercise: &domain.Exercise{
				Name:          "Squat",
				Description:   stringPtr("Back squat"),
				Aliases:       []string{"Back Squat", "BS"},
				MuscleGroups:  []string{"quadriceps", "glutes"},
				LoadIncrement: float64Ptr(2.5),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.exercise)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.exercise.ID == uuid.Nil {
					t.Error("Create() did not set ID")
				}
				if tt.exercise.CreatedAt.IsZero() {
					t.Error("Create() did not set CreatedAt")
				}
				if tt.exercise.UpdatedAt.IsZero() {
					t.Error("Create() did not set UpdatedAt")
				}
			}
		})
	}
}

func TestExerciseRepository_GetByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewExerciseRepository(pool)

	// Create an exercise first
	exercise := &domain.Exercise{
		Name:        "Test Exercise",
		Description: stringPtr("Test description"),
	}
	if err := repo.Create(ctx, exercise); err != nil {
		t.Fatalf("Failed to create exercise: %v", err)
	}

	// Test getting existing exercise
	got, err := repo.GetByID(ctx, exercise.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.ID != exercise.ID {
		t.Errorf("GetByID() ID = %v, want %v", got.ID, exercise.ID)
	}
	if got.Name != exercise.Name {
		t.Errorf("GetByID() Name = %v, want %v", got.Name, exercise.Name)
	}

	// Test getting non-existent exercise
	_, err = repo.GetByID(ctx, uuid.New())
	if err != domain.ErrNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestExerciseRepository_Update(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewExerciseRepository(pool)

	// Create an exercise
	exercise := &domain.Exercise{
		Name: "Original Name",
	}
	if err := repo.Create(ctx, exercise); err != nil {
		t.Fatalf("Failed to create exercise: %v", err)
	}

	originalUpdatedAt := exercise.UpdatedAt

	// Update the exercise
	exercise.Name = "Updated Name"
	exercise.Description = stringPtr("Updated description")
	if err := repo.Update(ctx, exercise); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if exercise.UpdatedAt.Equal(originalUpdatedAt) {
		t.Error("Update() did not update UpdatedAt")
	}

	// Verify update
	got, err := repo.GetByID(ctx, exercise.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.Name != "Updated Name" {
		t.Errorf("GetByID() Name = %v, want Updated Name", got.Name)
	}
}

func TestExerciseRepository_Delete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewExerciseRepository(pool)

	// Create an exercise
	exercise := &domain.Exercise{
		Name: "To Delete",
	}
	if err := repo.Create(ctx, exercise); err != nil {
		t.Fatalf("Failed to create exercise: %v", err)
	}

	// Delete the exercise
	if err := repo.Delete(ctx, exercise.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, exercise.ID)
	if err != domain.ErrNotFound {
		t.Errorf("GetByID() after Delete() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestExerciseRepository_List(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewExerciseRepository(pool)

	// Create multiple exercises
	for i := 0; i < 5; i++ {
		exercise := &domain.Exercise{
			Name: "Exercise " + string(rune('A'+i)),
		}
		if err := repo.Create(ctx, exercise); err != nil {
			t.Fatalf("Failed to create exercise: %v", err)
		}
	}

	// List exercises
	exercises, cursor, hasMore, err := repo.List(ctx, 3, "")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(exercises) != 3 {
		t.Errorf("List() returned %d exercises, want 3", len(exercises))
	}
	if !hasMore {
		t.Error("List() hasMore = false, want true")
	}
	if cursor == "" {
		t.Error("List() cursor is empty, want non-empty")
	}

	// List with cursor
	exercises2, _, hasMore2, err := repo.List(ctx, 3, cursor)
	if err != nil {
		t.Fatalf("List() with cursor error = %v", err)
	}

	// Should have remaining exercises (5 total - 3 already retrieved = 2 remaining)
	if len(exercises2) < 1 {
		t.Errorf("List() with cursor returned %d exercises, want at least 1", len(exercises2))
	}
	// hasMore2 should be false if we got all remaining exercises
	if len(exercises2) == 2 && hasMore2 {
		t.Error("List() with cursor hasMore = true, want false when all remaining items retrieved")
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
