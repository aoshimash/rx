package postgres

import (
	"context"
	"testing"

	"github.com/aoshimash/optel-training/api/internal/domain"
	"github.com/google/uuid"
)

func TestProgramRepository_Create(t *testing.T) {
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
			name: "create program with nodes",
			program: &domain.Program{
				Name:        "Program with Nodes",
				Description: stringPtr("Test description"),
				RootNodes: []domain.ProgramNode{
					{
						Name:     "Day 1",
						NodeType: "day",
						Order:    0,
						Children: []domain.ProgramNode{
							{
								Name:     "Exercise 1",
								NodeType: "exercise",
								Order:    0,
							},
						},
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
		RootNodes: []domain.ProgramNode{
			{
				Name:     "Day 1",
				NodeType: "day",
				Order:    0,
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
	if len(got.RootNodes) != 1 {
		t.Errorf("GetByID() RootNodes length = %v, want 1", len(got.RootNodes))
	}

	// Test getting non-existent program
	_, err = repo.GetByID(ctx, uuid.New())
	if err != domain.ErrNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, domain.ErrNotFound)
	}
}

func TestProgramRepository_Update(t *testing.T) {
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

func TestProgramRepository_GetByID_NestedTree(t *testing.T) {
	ctx := context.Background()
	pool, cleanup, err := setupTestDB(ctx)
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer cleanup()

	repo := NewProgramRepository(pool)

	// Create a program with deeply nested structure: Program -> Week -> Day -> Exercise
	program := &domain.Program{
		Name:        "Deep Nested Program",
		Description: stringPtr("Test nested tree structure"),
		RootNodes: []domain.ProgramNode{
			{
				Name:     "Week 1",
				NodeType: "week",
				Order:    0,
				Children: []domain.ProgramNode{
					{
						Name:     "Day 1",
						NodeType: "day",
						Order:    0,
						Children: []domain.ProgramNode{
							{
								Name:     "Exercise 1",
								NodeType: "exercise",
								Order:    0,
							},
							{
								Name:     "Exercise 2",
								NodeType: "exercise",
								Order:    1,
							},
						},
					},
					{
						Name:     "Day 2",
						NodeType: "day",
						Order:    1,
						Children: []domain.ProgramNode{
							{
								Name:     "Exercise 3",
								NodeType: "exercise",
								Order:    0,
							},
						},
					},
				},
			},
			{
				Name:     "Week 2",
				NodeType: "week",
				Order:    1,
				Children: []domain.ProgramNode{
					{
						Name:     "Day 3",
						NodeType: "day",
						Order:    0,
					},
				},
			},
		},
	}

	if err := repo.Create(ctx, program); err != nil {
		t.Fatalf("Failed to create program: %v", err)
	}

	// Retrieve and verify nested structure
	got, err := repo.GetByID(ctx, program.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	// Verify root level
	if len(got.RootNodes) != 2 {
		t.Fatalf("Expected 2 root nodes (weeks), got %d", len(got.RootNodes))
	}

	// Find Week 1 (may be in any order due to UUID sorting)
	var week1, week2 *domain.ProgramNode
	for i := range got.RootNodes {
		switch got.RootNodes[i].Name {
		case "Week 1":
			week1 = &got.RootNodes[i]
		case "Week 2":
			week2 = &got.RootNodes[i]
		}
	}

	if week1 == nil {
		t.Fatal("Week 1 not found in root nodes")
	}
	if week2 == nil {
		t.Fatal("Week 2 not found in root nodes")
	}

	// Verify Week 1 has 2 days
	if len(week1.Children) != 2 {
		t.Errorf("Week 1 should have 2 children (days), got %d", len(week1.Children))
	}

	// Find Day 1 in Week 1
	var day1 *domain.ProgramNode
	for i := range week1.Children {
		if week1.Children[i].Name == "Day 1" {
			day1 = &week1.Children[i]
			break
		}
	}

	if day1 == nil {
		t.Fatal("Day 1 not found in Week 1")
	}

	// Verify Day 1 has 2 exercises (the critical test for nested children)
	if len(day1.Children) != 2 {
		t.Errorf("Day 1 should have 2 children (exercises), got %d. Nested children may be lost!", len(day1.Children))
	}

	// Verify Week 2 has 1 day
	if len(week2.Children) != 1 {
		t.Errorf("Week 2 should have 1 child (day), got %d", len(week2.Children))
	}
}
