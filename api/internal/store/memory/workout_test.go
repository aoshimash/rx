package memory

import (
	"context"
	"testing"
	"time"

	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

func TestWorkoutRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		workout *domain.Workout
		wantErr bool
	}{
		{
			name: "create valid workout",
			workout: &domain.Workout{
				Timestamp: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "create workout with entries",
			workout: &domain.Workout{
				Timestamp: time.Now(),
				Entries: []domain.WorkoutEntry{
					{
						Order:      1,
						ExerciseID: uuid.New(),
						EntryType:  stringPtr("work"),
						Sets:       3,
						Reps:       10,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewWorkoutRepository()
			ctx := context.Background()

			err := repo.Create(ctx, tt.workout)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.workout.ID == uuid.Nil {
					t.Error("Create() did not generate ID")
				}
				if tt.workout.CreatedAt.IsZero() {
					t.Error("Create() did not set CreatedAt")
				}
				if len(tt.workout.Entries) > 0 {
					if tt.workout.Entries[0].ID == uuid.Nil {
						t.Error("Create() did not generate Entry ID")
					}
					if tt.workout.Entries[0].WorkoutID != tt.workout.ID {
						t.Error("Create() did not set Entry WorkoutID")
					}
				}
			}
		})
	}
}

func TestWorkoutRepository_GetByID(t *testing.T) {
	repo := NewWorkoutRepository().(*workoutStore)
	ctx := context.Background()

	id := uuid.New()
	exerciseID := uuid.New()
	workout := &domain.Workout{
		ID:        id,
		Timestamp: time.Now(),
		Entries: []domain.WorkoutEntry{
			{
				ID:         uuid.New(),
				WorkoutID:  id,
				Order:      1,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("work"),
				Sets:       3,
				Reps:       10,
				LoadKg:     100.0,
				RPE:        8,
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.workouts[id] = workout

	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.ID != id {
		t.Errorf("GetByID() got ID = %v, want %v", got.ID, id)
	}
	if len(got.Entries) != 1 {
		t.Errorf("GetByID() got %d entries, want 1", len(got.Entries))
	}
	if got.Entries[0].ExerciseID != exerciseID {
		t.Errorf("GetByID() got Entry ExerciseID = %v, want %v", got.Entries[0].ExerciseID, exerciseID)
	}

	// Verify it's a copy
	originalName := got.Entries[0].ExerciseID
	got.Entries[0].ExerciseID = uuid.New()
	recheck, _ := repo.GetByID(ctx, id)
	if recheck.Entries[0].ExerciseID != originalName {
		t.Error("GetByID() did not return a copy")
	}
}

func TestWorkoutRepository_ListByExerciseID(t *testing.T) {
	repo := NewWorkoutRepository().(*workoutStore)
	ctx := context.Background()

	exerciseID := uuid.New()
	workout1 := &domain.Workout{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		Entries: []domain.WorkoutEntry{
			{ID: uuid.New(), WorkoutID: uuid.New(), ExerciseID: exerciseID},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	workout2 := &domain.Workout{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		Entries: []domain.WorkoutEntry{
			{ID: uuid.New(), WorkoutID: uuid.New(), ExerciseID: uuid.New()},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.workouts[workout1.ID] = workout1
	repo.workouts[workout2.ID] = workout2

	workouts, err := repo.ListByExerciseID(ctx, exerciseID)
	if err != nil {
		t.Fatalf("ListByExerciseID() error = %v", err)
	}

	if len(workouts) != 1 {
		t.Errorf("ListByExerciseID() got %d workouts, want 1", len(workouts))
	}
	if workouts[0].ID != workout1.ID {
		t.Errorf("ListByExerciseID() got workout ID = %v, want %v", workouts[0].ID, workout1.ID)
	}
}

func TestWorkoutRepository_ListByTimestampRange(t *testing.T) {
	repo := NewWorkoutRepository().(*workoutStore)
	ctx := context.Background()

	now := time.Now()
	workout1 := &domain.Workout{
		ID:        uuid.New(),
		Timestamp: now.Add(-2 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	workout2 := &domain.Workout{
		ID:        uuid.New(),
		Timestamp: now.Add(-1 * time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	workout3 := &domain.Workout{
		ID:        uuid.New(),
		Timestamp: now,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.workouts[workout1.ID] = workout1
	repo.workouts[workout2.ID] = workout2
	repo.workouts[workout3.ID] = workout3

	from := now.Add(-90 * time.Minute)
	to := now.Add(-30 * time.Minute)

	workouts, _, _, err := repo.ListByTimestampRange(ctx, &from, &to, 10, "")
	if err != nil {
		t.Fatalf("ListByTimestampRange() error = %v", err)
	}

	if len(workouts) != 1 {
		t.Errorf("ListByTimestampRange() got %d workouts, want 1", len(workouts))
	}
	if workouts[0].ID != workout2.ID {
		t.Errorf("ListByTimestampRange() got workout ID = %v, want %v", workouts[0].ID, workout2.ID)
	}
}
