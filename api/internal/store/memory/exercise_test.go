package memory

import (
	"context"
	"testing"
	"time"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/google/uuid"
)

func TestExerciseRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		exercise *domain.Exercise
		wantErr bool
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
			repo := NewExerciseRepository()
			ctx := context.Background()

			err := repo.Create(ctx, tt.exercise)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.exercise.ID == uuid.Nil {
					t.Error("Create() did not generate ID")
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
	tests := []struct {
		name    string
		setup   func(*exerciseStore)
		id      uuid.UUID
		wantErr bool
		wantID  uuid.UUID
	}{
		{
			name: "get existing exercise",
			setup: func(s *exerciseStore) {
				id := uuid.New()
				s.exercises[id] = &domain.Exercise{
					ID:        id,
					Name:      "Bench Press",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
			},
			wantErr: false,
		},
		{
			name:    "get non-existent exercise",
			setup:   func(s *exerciseStore) {},
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewExerciseRepository().(*exerciseStore)
			ctx := context.Background()

			tt.setup(repo)

			var id uuid.UUID
			if tt.id == uuid.Nil {
				// Get first exercise ID from store
				for k := range repo.exercises {
					id = k
					break
				}
			} else {
				id = tt.id
			}

			exercise, err := repo.GetByID(ctx, id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if exercise == nil {
					t.Error("GetByID() returned nil exercise")
					return
				}
				if exercise.ID != id {
					t.Errorf("GetByID() got ID = %v, want %v", exercise.ID, id)
				}
				// Verify it's a copy (modifying shouldn't affect store)
				originalName := exercise.Name
				exercise.Name = "Modified"
				recheck, _ := repo.GetByID(ctx, id)
				if recheck.Name != originalName {
					t.Error("GetByID() did not return a copy")
				}
			} else {
				if err != domain.ErrNotFound {
					t.Errorf("GetByID() error = %v, want %v", err, domain.ErrNotFound)
				}
			}
		})
	}
}

func TestExerciseRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*exerciseStore) uuid.UUID
		update  func(*domain.Exercise)
		wantErr bool
	}{
		{
			name: "update existing exercise",
			setup: func(s *exerciseStore) uuid.UUID {
				id := uuid.New()
				s.exercises[id] = &domain.Exercise{
					ID:        id,
					Name:      "Bench Press",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				return id
			},
			update: func(e *domain.Exercise) {
				e.Name = "Modified Bench Press"
			},
			wantErr: false,
		},
		{
			name: "update non-existent exercise",
			setup: func(s *exerciseStore) uuid.UUID {
				return uuid.New()
			},
			update: func(e *domain.Exercise) {
				e.Name = "New Exercise"
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewExerciseRepository().(*exerciseStore)
			ctx := context.Background()

			id := tt.setup(repo)
			exercise := &domain.Exercise{
				ID:        id,
				Name:      "Original",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			tt.update(exercise)

			err := repo.Update(ctx, exercise)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Get original UpdatedAt before update
				original, _ := repo.GetByID(ctx, id)
				originalUpdatedAt := original.UpdatedAt

				// Wait a bit to ensure timestamp difference
				time.Sleep(10 * time.Millisecond)

				// Update again to trigger UpdatedAt change
				err = repo.Update(ctx, exercise)
				if err != nil {
					t.Fatalf("Update() failed on second call: %v", err)
				}

				updated, _ := repo.GetByID(ctx, id)
				if updated.Name != exercise.Name {
					t.Errorf("Update() did not update exercise, got Name = %v, want %v", updated.Name, exercise.Name)
				}
				if !updated.UpdatedAt.After(originalUpdatedAt) {
					t.Errorf("Update() did not update UpdatedAt timestamp: got %v, want after %v", updated.UpdatedAt, originalUpdatedAt)
				}
			}
		})
	}
}

func TestExerciseRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*exerciseStore) uuid.UUID
		wantErr bool
	}{
		{
			name: "delete existing exercise",
			setup: func(s *exerciseStore) uuid.UUID {
				id := uuid.New()
				s.exercises[id] = &domain.Exercise{
					ID:        id,
					Name:      "Bench Press",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				return id
			},
			wantErr: false,
		},
		{
			name: "delete non-existent exercise",
			setup: func(s *exerciseStore) uuid.UUID {
				return uuid.New()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewExerciseRepository().(*exerciseStore)
			ctx := context.Background()

			id := tt.setup(repo)

			err := repo.Delete(ctx, id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				_, err := repo.GetByID(ctx, id)
				if err != domain.ErrNotFound {
					t.Error("Delete() did not remove exercise")
				}
			}
		})
	}
}

func TestExerciseRepository_List(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*exerciseStore)
		limit     int
		after     string
		wantCount int
		wantMore  bool
		wantErr   bool
	}{
		{
			name: "list all exercises with limit",
			setup: func(s *exerciseStore) {
				for i := 0; i < 5; i++ {
					id := uuid.New()
					s.exercises[id] = &domain.Exercise{
						ID:        id,
						Name:      "Exercise",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
				}
			},
			limit:     3,
			after:     "",
			wantCount: 3,
			wantMore:  true,
			wantErr:   false,
		},
		{
			name: "list with pagination cursor",
			setup: func(s *exerciseStore) {
				for i := 0; i < 5; i++ {
					id := uuid.New()
					s.exercises[id] = &domain.Exercise{
						ID:        id,
						Name:      "Exercise",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}
				}
			},
			limit:     2,
			after:     "", // Will be set after first call
			wantCount: 2,
			wantMore:  true,
			wantErr:   false,
		},
		{
			name: "list empty store",
			setup: func(s *exerciseStore) {
				// No exercises
			},
			limit:     10,
			after:     "",
			wantCount: 0,
			wantMore:  false,
			wantErr:   false,
		},
		{
			name: "list with invalid cursor",
			setup: func(s *exerciseStore) {
				// No exercises needed
			},
			limit:     10,
			after:     "invalid-cursor",
			wantCount: 0,
			wantMore:  false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewExerciseRepository().(*exerciseStore)
			ctx := context.Background()

			tt.setup(repo)

			after := tt.after
			if tt.name == "list with pagination cursor" {
				// Get first page to get cursor
				first, cursor, _, _ := repo.List(ctx, 2, "")
				if len(first) == 0 {
					t.Fatal("Expected at least one exercise for pagination test")
				}
				after = cursor
			}

			exercises, nextCursor, hasMore, err := repo.List(ctx, tt.limit, after)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(exercises) != tt.wantCount {
					t.Errorf("List() got %d exercises, want %d", len(exercises), tt.wantCount)
				}
				if hasMore != tt.wantMore {
					t.Errorf("List() got hasMore = %v, want %v", hasMore, tt.wantMore)
				}
				if tt.wantMore && nextCursor == "" {
					t.Error("List() should return next cursor when hasMore is true")
				}
				if !tt.wantMore && nextCursor != "" {
					t.Error("List() should not return next cursor when hasMore is false")
				}
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
