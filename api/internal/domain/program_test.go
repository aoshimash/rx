package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateProgramEntry(t *testing.T) {
	programID := uuid.New()
	exerciseID := uuid.New()

	tests := []struct {
		name    string
		entry   *ProgramEntry
		wantErr bool
		errCode string
	}{
		{
			name:    "nil entry",
			entry:   nil,
			wantErr: true,
		},
		{
			name: "valid entry with minimal fields",
			entry: &ProgramEntry{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Bench Press",
				Order:     0,
			},
			wantErr: false,
		},
		{
			name: "valid entry with metadata",
			entry: &ProgramEntry{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Squat",
				Order:     1,
				Metadata:  []byte(`{"week": 1, "day": 2}`),
			},
			wantErr: false,
		},
		{
			name: "valid entry with all prescription fields",
			entry: &ProgramEntry{
				ID:                 uuid.New(),
				ProgramID:          programID,
				Name:               "Bench Press",
				Order:              0,
				ExerciseID:         &exerciseID,
				TargetSets:         intPtr(5),
				TargetReps:         intPtr(5),
				TargetRPE:          intPtr(7),
				Percent1RM:         float64Ptr(0.75),
				PlannedRestSeconds: intPtr(180),
				MuscleGroups:       []string{"chest", "triceps"},
				Notes:              stringPtr("Focus on form"),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			entry: &ProgramEntry{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "",
				Order:     0,
			},
			wantErr: true,
			errCode: ErrCodeMissingRequiredField,
		},
		{
			name: "name too long",
			entry: &ProgramEntry{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      stringWithLength(201),
				Order:     0,
			},
			wantErr: true,
		},
		{
			name: "order negative",
			entry: &ProgramEntry{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Exercise",
				Order:     -1,
			},
			wantErr: true,
		},
		{
			name: "target_sets zero",
			entry: &ProgramEntry{
				ID:         uuid.New(),
				ProgramID:  programID,
				Name:       "Exercise",
				Order:      0,
				TargetSets: intPtr(0),
			},
			wantErr: true,
		},
		{
			name: "target_reps zero",
			entry: &ProgramEntry{
				ID:         uuid.New(),
				ProgramID:  programID,
				Name:       "Exercise",
				Order:      0,
				TargetReps: intPtr(0),
			},
			wantErr: true,
		},
		{
			name: "target_rpe invalid",
			entry: &ProgramEntry{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Exercise",
				Order:     0,
				TargetRPE: intPtr(11),
			},
			wantErr: true,
			errCode: ErrCodeInvalidRPE,
		},
		{
			name: "percent_1rm too low",
			entry: &ProgramEntry{
				ID:         uuid.New(),
				ProgramID:  programID,
				Name:       "Exercise",
				Order:      0,
				Percent1RM: float64Ptr(-0.1),
			},
			wantErr: true,
		},
		{
			name: "percent_1rm too high",
			entry: &ProgramEntry{
				ID:         uuid.New(),
				ProgramID:  programID,
				Name:       "Exercise",
				Order:      0,
				Percent1RM: float64Ptr(1.1),
			},
			wantErr: true,
		},
		{
			name: "planned_rest_seconds negative",
			entry: &ProgramEntry{
				ID:                 uuid.New(),
				ProgramID:          programID,
				Name:               "Exercise",
				Order:              0,
				PlannedRestSeconds: intPtr(-1),
			},
			wantErr: true,
		},
		{
			name: "notes too long",
			entry: &ProgramEntry{
				ID:        uuid.New(),
				ProgramID: programID,
				Name:      "Exercise",
				Order:     0,
				Notes:     stringPtr(stringWithLength(2001)),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProgramEntry(tt.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProgramEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				if domainErr, ok := err.(*DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("ValidateProgramEntry() error code = %v, want %v", domainErr.Code, tt.errCode)
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
			name: "valid program with empty entries",
			program: &Program{
				ID:        uuid.New(),
				Name:      "Program",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Entries:   []ProgramEntry{},
			},
			wantErr: false,
		},
		{
			name: "valid program with entries",
			program: &Program{
				ID:        uuid.New(),
				Name:      "Program",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Entries: []ProgramEntry{
					{
						ID:         uuid.New(),
						ProgramID:  programID,
						Name:       "Squat",
						Order:      0,
						Metadata:   []byte(`{"week": 1, "day": 1}`),
						ExerciseID: &exerciseID,
						TargetSets: intPtr(5),
						TargetReps: intPtr(5),
						TargetRPE:  intPtr(7),
					},
					{
						ID:        uuid.New(),
						ProgramID: programID,
						Name:      "Bench Press",
						Order:     1,
						Metadata:  []byte(`{"week": 1, "day": 1}`),
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
			name: "invalid entry",
			program: &Program{
				ID:        uuid.New(),
				Name:      "Program",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Entries: []ProgramEntry{
					{
						ID:        uuid.New(),
						ProgramID: programID,
						Name:      "", // Invalid
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
