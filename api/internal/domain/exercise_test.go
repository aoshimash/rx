package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateExercise(t *testing.T) {
	tests := []struct {
		name    string
		exercise *Exercise
		wantErr bool
		errCode string
	}{
		{
			name:    "nil exercise",
			exercise: nil,
			wantErr: true,
		},
		{
			name: "valid exercise with minimal fields",
			exercise: &Exercise{
				ID:        uuid.New(),
				Name:      "Bench Press",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid exercise with all fields",
			exercise: &Exercise{
				ID:            uuid.New(),
				Name:          "Squat",
				Description:   stringPtr("Back squat"),
				Aliases:       []string{"Back Squat", "BS"},
				MuscleGroups:  []string{"quadriceps", "glutes"},
				LoadIncrement: float64Ptr(2.5),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			exercise: &Exercise{
				ID:        uuid.New(),
				Name:      "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "name too long",
			exercise: &Exercise{
				ID:        uuid.New(),
				Name:      stringWithLength(201),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "description too long",
			exercise: &Exercise{
				ID:          uuid.New(),
				Name:        "Exercise",
				Description: stringPtr(stringWithLength(2001)),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			wantErr: true,
		},
		{
			name: "load_increment zero",
			exercise: &Exercise{
				ID:            uuid.New(),
				Name:          "Exercise",
				LoadIncrement: float64Ptr(0),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			wantErr: true,
		},
		{
			name: "load_increment negative",
			exercise: &Exercise{
				ID:            uuid.New(),
				Name:          "Exercise",
				LoadIncrement: float64Ptr(-1.0),
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExercise(tt.exercise)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExercise() error = %v, wantErr %v", err, tt.wantErr)
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

func stringWithLength(n int) string {
	result := make([]byte, n)
	for i := range result {
		result[i] = 'a'
	}
	return string(result)
}
