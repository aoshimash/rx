package memory

import (
	"context"
	"testing"
	"time"

	"github.com/aoshimash/optel-workout/api/internal/domain"
	"github.com/google/uuid"
)

func TestProgramRepository_Create(t *testing.T) {
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
			name: "create program with root nodes",
			program: &domain.Program{
				Name: "Test Program",
				RootNodes: []domain.ProgramNode{
					{
						Name: "Week 1",
						Children: []domain.ProgramNode{
							{Name: "Day 1"},
							{Name: "Day 2"},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewProgramRepository()
			ctx := context.Background()

			err := repo.Create(ctx, tt.program)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.program.ID == uuid.Nil {
					t.Error("Create() did not generate ID")
				}
				if tt.program.CreatedAt.IsZero() {
					t.Error("Create() did not set CreatedAt")
				}
				if len(tt.program.RootNodes) > 0 {
					if tt.program.RootNodes[0].ID == uuid.Nil {
						t.Error("Create() did not generate RootNode ID")
					}
					if len(tt.program.RootNodes[0].Children) > 0 {
						if tt.program.RootNodes[0].Children[0].ID == uuid.Nil {
							t.Error("Create() did not generate Child Node ID")
						}
					}
				}
			}
		})
	}
}

func TestProgramRepository_GetByID(t *testing.T) {
	repo := NewProgramRepository().(*programStore)
	ctx := context.Background()

	id := uuid.New()
	program := &domain.Program{
		ID:   id,
		Name: "Test Program",
		RootNodes: []domain.ProgramNode{
			{
				ID:   uuid.New(),
				Name: "Week 1",
				Children: []domain.ProgramNode{
					{ID: uuid.New(), Name: "Day 1"},
				},
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.programs[id] = program

	got, err := repo.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if got.ID != id {
		t.Errorf("GetByID() got ID = %v, want %v", got.ID, id)
	}
	if len(got.RootNodes) != 1 {
		t.Errorf("GetByID() got %d root nodes, want 1", len(got.RootNodes))
	}
	if len(got.RootNodes[0].Children) != 1 {
		t.Errorf("GetByID() got %d children, want 1", len(got.RootNodes[0].Children))
	}

	// Verify it's a deep copy
	originalName := got.RootNodes[0].Name
	got.RootNodes[0].Name = "Modified"
	recheck, _ := repo.GetByID(ctx, id)
	if recheck.RootNodes[0].Name != originalName {
		t.Error("GetByID() did not return a deep copy")
	}
}

func TestProgramRepository_Update(t *testing.T) {
	repo := NewProgramRepository().(*programStore)
	ctx := context.Background()

	id := uuid.New()
	program := &domain.Program{
		ID:        id,
		Name:      "Original Program",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.programs[id] = program

	program.Name = "Updated Program"
	err := repo.Update(ctx, program)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	updated, _ := repo.GetByID(ctx, id)
	if updated.Name != "Updated Program" {
		t.Errorf("Update() did not update program, got Name = %v, want Updated Program", updated.Name)
	}
}

func TestProgramRepository_Delete(t *testing.T) {
	repo := NewProgramRepository().(*programStore)
	ctx := context.Background()

	id := uuid.New()
	program := &domain.Program{
		ID:        id,
		Name:      "Test Program",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.programs[id] = program

	err := repo.Delete(ctx, id)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = repo.GetByID(ctx, id)
	if err != domain.ErrNotFound {
		t.Error("Delete() did not remove program")
	}
}
