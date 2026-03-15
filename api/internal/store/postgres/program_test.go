package postgres

import (
	"context"
	"testing"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/google/uuid"
)

func TestProgramRepository_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewProgramRepository(pool)

	tests := []struct {
		name    string
		program *domain.Program
		wantErr bool
	}{
		{
			name: "create valid program",
			program: &domain.Program{
				Name: "Test Program",
			},
			wantErr: false,
		},
		{
			name: "create program with entries",
			program: &domain.Program{
				Name:        "Program with Entries",
				Description: stringPtr("Test description"),
				Entries: []domain.ProgramEntry{
					{
						Name:     "Squat",
						Order:    0,
						Metadata: []byte(`{"week": "Week 1", "day": "Day 1"}`),
					},
					{
						Name:  "Bench Press",
						Order: 1,
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(ctx, tt.program)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.program.ID == uuid.Nil {
					t.Error("Create() did not set ID")
				}
				if tt.program.CreatedAt.IsZero() {
					t.Error("Create() did not set CreatedAt")
				}
			}
		})
	}
}

func TestProgramRepository_GetByID(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewProgramRepository(pool)

	// Create a program
	program := &domain.Program{
		Name:        "Test Program",
		Description: stringPtr("Test description"),
		Entries: []domain.ProgramEntry{
			{
				Name:     "Squat",
				Order:    0,
				Metadata: []byte(`{"week": "Week 1", "day": "Day 1"}`),
			},
		},
	}
	if err := repo.Create(ctx, program); err != nil {
		t.Fatalf("Failed to create program: %v", err)
	}

	// Test getting existing program
	got, err := repo.GetByID(ctx, program.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.ID != program.ID {
		t.Errorf("GetByID() ID = %v, want %v", got.ID, program.ID)
	}
	if got.Name != program.Name {
		t.Errorf("GetByID() Name = %v, want %v", got.Name, program.Name)
	}
	if len(got.Entries) != 1 {
		t.Errorf("GetByID() Entries length = %v, want 1", len(got.Entries))
	}

	// Test getting non-existent program
	_, err = repo.GetByID(ctx, uuid.New())
	if err != domain.ErrNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestProgramRepository_Update(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewProgramRepository(pool)

	// Create a program
	program := &domain.Program{
		Name: "Original Name",
	}
	if err := repo.Create(ctx, program); err != nil {
		t.Fatalf("Failed to create program: %v", err)
	}

	// Update the program
	program.Name = "Updated Name"
	program.Description = stringPtr("Updated description")
	if err := repo.Update(ctx, program); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	// Verify update
	got, err := repo.GetByID(ctx, program.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.Name != "Updated Name" {
		t.Errorf("GetByID() Name = %v, want Updated Name", got.Name)
	}
}

func TestProgramRepository_Delete(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewProgramRepository(pool)

	// Create a program
	program := &domain.Program{
		Name: "To Delete",
	}
	if err := repo.Create(ctx, program); err != nil {
		t.Fatalf("Failed to create program: %v", err)
	}

	// Delete the program
	if err := repo.Delete(ctx, program.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify deletion
	_, err = repo.GetByID(ctx, program.ID)
	if err != domain.ErrNotFound {
		t.Errorf("GetByID() after Delete() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestProgramRepository_List(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewProgramRepository(pool)

	// Create multiple programs
	for i := 0; i < 3; i++ {
		program := &domain.Program{
			Name: "Program " + string(rune('A'+i)),
		}
		if err := repo.Create(ctx, program); err != nil {
			t.Fatalf("Failed to create program: %v", err)
		}
	}

	// List programs
	programs, _, hasMore, err := repo.List(ctx, 2, "")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(programs) != 2 {
		t.Errorf("List() returned %d programs, want 2", len(programs))
	}
	if !hasMore {
		t.Error("List() hasMore = false, want true")
	}
}

func TestProgramRepository_List_WithEntries(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewProgramRepository(pool)

	// Create a program with entries
	programWithEntries := &domain.Program{
		Name: "Program With Entries",
		Entries: []domain.ProgramEntry{
			{
				Name:     "Squat",
				Order:    0,
				Metadata: []byte(`{"week": "Week 1", "day": "Day 1"}`),
			},
			{
				Name:  "Bench Press",
				Order: 1,
			},
		},
	}
	if err := repo.Create(ctx, programWithEntries); err != nil {
		t.Fatalf("Failed to create program with entries: %v", err)
	}

	// Create a program without entries
	programWithoutEntries := &domain.Program{
		Name: "Program Without Entries",
	}
	if err := repo.Create(ctx, programWithoutEntries); err != nil {
		t.Fatalf("Failed to create program without entries: %v", err)
	}

	// List programs
	programs, _, _, err := repo.List(ctx, 10, "")
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(programs) != 2 {
		t.Fatalf("List() returned %d programs, want 2", len(programs))
	}

	// Find program with entries and verify entries are loaded
	var foundProgramWithEntries *domain.Program
	for _, p := range programs {
		if p.Name == "Program With Entries" {
			foundProgramWithEntries = p
			break
		}
	}

	if foundProgramWithEntries == nil {
		t.Fatal("Program With Entries not found in list")
	}

	if len(foundProgramWithEntries.Entries) != 2 {
		t.Errorf("List() should load Entries, got %d, want 2", len(foundProgramWithEntries.Entries))
	}
}

func TestProgramRepository_GetByID_WithMetadata(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewProgramRepository(pool)

	// Create a program with entries containing metadata
	program := &domain.Program{
		Name:        "Flat Program",
		Description: stringPtr("Test flat structure"),
		Entries: []domain.ProgramEntry{
			{
				Name:     "Squat",
				Order:    0,
				Metadata: []byte(`{"week": "Week 1", "day": "Day 1"}`),
			},
			{
				Name:     "Bench Press",
				Order:    1,
				Metadata: []byte(`{"week": "Week 1", "day": "Day 1"}`),
			},
			{
				Name:     "Deadlift",
				Order:    2,
				Metadata: []byte(`{"week": "Week 1", "day": "Day 2"}`),
			},
		},
	}

	if err := repo.Create(ctx, program); err != nil {
		t.Fatalf("Failed to create program: %v", err)
	}

	// Retrieve and verify flat structure
	got, err := repo.GetByID(ctx, program.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if len(got.Entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(got.Entries))
	}

	// Verify metadata is preserved
	if string(got.Entries[0].Metadata) == "" {
		t.Error("Entry 0 should have metadata")
	}
}
