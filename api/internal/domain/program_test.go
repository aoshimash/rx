package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateProgramNode(t *testing.T) {
	programID := uuid.New()
	parentID := uuid.New()
	exerciseID := uuid.New()

	tests := []struct {
		name    string
		node    *ProgramNode
		wantErr bool
		errCode string
	}{
		{
			name:    "nil node",
			node:    nil,
			wantErr: true,
		},
		{
			name: "valid node with minimal fields",
			node: &ProgramNode{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Week 1",
				NodeType:  "week",
				Order:     0,
			},
			wantErr: false,
		},
		{
			name: "valid node with parent",
			node: &ProgramNode{
				ID:        uuid.New(),
				ProgramID: programID,
				ParentID:  &parentID,
				Name:      "Day 1",
				NodeType:  "day",
				Order:     0,
			},
			wantErr: false,
		},
		{
			name: "valid leaf node with prescription fields",
			node: &ProgramNode{
				ID:                uuid.New(),
				ProgramID:         programID,
				ParentID:          &parentID,
				Name:              "Bench Press",
				NodeType:          "exercise",
				Order:             0,
				ExerciseID:        &exerciseID,
				TargetSets:        intPtr(5),
				TargetReps:        intPtr(5),
				TargetRPE:         intPtr(7),
				Percent1RM:         float64Ptr(0.75),
				PlannedRestSeconds: intPtr(180),
				MuscleGroups:      []string{"chest", "triceps"},
				Notes:             stringPtr("Focus on form"),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			node: &ProgramNode{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "",
				NodeType:  "week",
				Order:     0,
			},
			wantErr: true,
			errCode: ErrCodeMissingRequiredField,
		},
		{
			name: "name too long",
			node: &ProgramNode{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      stringWithLength(201),
				NodeType:  "week",
				Order:     0,
			},
			wantErr: true,
		},
		{
			name: "empty node_type",
			node: &ProgramNode{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Node",
				NodeType:  "",
				Order:     0,
			},
			wantErr: true,
			errCode: ErrCodeMissingRequiredField,
		},
		{
			name: "order negative",
			node: &ProgramNode{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Node",
				NodeType:  "week",
				Order:     -1,
			},
			wantErr: true,
		},
		{
			name: "target_sets zero",
			node: &ProgramNode{
				ID:         uuid.New(),
				ProgramID:  programID,
				Name:       "Exercise",
				NodeType:   "exercise",
				Order:      0,
				TargetSets: intPtr(0),
			},
			wantErr: true,
		},
		{
			name: "target_reps zero",
			node: &ProgramNode{
				ID:         uuid.New(),
				ProgramID:  programID,
				Name:       "Exercise",
				NodeType:   "exercise",
				Order:      0,
				TargetReps: intPtr(0),
			},
			wantErr: true,
		},
		{
			name: "target_rpe invalid",
			node: &ProgramNode{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Exercise",
				NodeType:  "exercise",
				Order:     0,
				TargetRPE: intPtr(11),
			},
			wantErr: true,
			errCode: ErrCodeInvalidRPE,
		},
		{
			name: "percent_1rm too low",
			node: &ProgramNode{
				ID:         uuid.New(),
				ProgramID:  programID,
				Name:       "Exercise",
				NodeType:   "exercise",
				Order:      0,
				Percent1RM: float64Ptr(-0.1),
			},
			wantErr: true,
		},
		{
			name: "percent_1rm too high",
			node: &ProgramNode{
				ID:         uuid.New(),
				ProgramID:  programID,
				Name:       "Exercise",
				NodeType:   "exercise",
				Order:      0,
				Percent1RM: float64Ptr(1.1),
			},
			wantErr: true,
		},
		{
			name: "planned_rest_seconds negative",
			node: &ProgramNode{
				ID:                uuid.New(),
				ProgramID:         programID,
				Name:              "Exercise",
				NodeType:          "exercise",
				Order:             0,
				PlannedRestSeconds: intPtr(-1),
			},
			wantErr: true,
		},
		{
			name: "notes too long",
			node: &ProgramNode{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Node",
				NodeType:  "week",
				Order:     0,
				Notes:     stringPtr(stringWithLength(2001)),
			},
			wantErr: true,
		},
		{
			name: "recursive children validation",
			node: &ProgramNode{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Cycle 1",
				NodeType:  "cycle",
				Order:     0,
				Children: []ProgramNode{
					{
						ID:        uuid.New(),
						ProgramID: programID,
						Name:      "", // Invalid: empty name
						NodeType:  "week",
						Order:     0,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProgramNode(tt.node)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProgramNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				if domainErr, ok := err.(*DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("ValidateProgramNode() error code = %v, want %v", domainErr.Code, tt.errCode)
					}
				}
			}
		})
	}
}

func TestValidateProgram(t *testing.T) {
	programID := uuid.New()
	exerciseID := uuid.New()

	tests := []struct {
		name    string
		program *Program
		wantErr bool
		errCode string
	}{
		{
			name:    "nil program",
			program: nil,
			wantErr: true,
		},
		{
			name: "valid program with minimal fields",
			program: &Program{
				ID:        uuid.New(),
				Name:      "Full Body Hypertrophy",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid program with empty root_nodes",
			program: &Program{
				ID:        uuid.New(),
				Name:      "Program",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				RootNodes: []ProgramNode{},
			},
			wantErr: false,
		},
		{
			name: "valid program with root nodes",
			program: &Program{
				ID:        uuid.New(),
				Name:      "Program",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				RootNodes: []ProgramNode{
					{
						ID:        uuid.New(),
						ProgramID: programID,
						Name:      "Cycle 1",
						NodeType:  "cycle",
						Order:     0,
						Children: []ProgramNode{
							{
								ID:        uuid.New(),
								ProgramID: programID,
								Name:      "Week 1",
								NodeType:  "week",
								Order:     0,
								Children: []ProgramNode{
									{
										ID:         uuid.New(),
										ProgramID:  programID,
										Name:       "Day 1",
										NodeType:   "day",
										Order:      0,
										ExerciseID: &exerciseID,
										TargetSets: intPtr(5),
										TargetReps: intPtr(5),
										TargetRPE:  intPtr(7),
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty name",
			program: &Program{
				ID:        uuid.New(),
				Name:      "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
			errCode: ErrCodeMissingRequiredField,
		},
		{
			name: "name too long",
			program: &Program{
				ID:        uuid.New(),
				Name:      stringWithLength(201),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "description too long",
			program: &Program{
				ID:          uuid.New(),
				Name:        "Program",
				Description: stringPtr(stringWithLength(2001)),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "invalid root node",
			program: &Program{
				ID:        uuid.New(),
				Name:      "Program",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				RootNodes: []ProgramNode{
					{
						ID:        uuid.New(),
						ProgramID: programID,
						Name:      "", // Invalid
						NodeType:  "week",
						Order:     0,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProgram(tt.program)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				if domainErr, ok := err.(*DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("ValidateProgram() error code = %v, want %v", domainErr.Code, tt.errCode)
					}
				}
			}
		})
	}
}
