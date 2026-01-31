package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateWorkoutEntry(t *testing.T) {
	now := time.Now()
	exerciseID := uuid.New()
	workoutID := uuid.New()

	tests := []struct {
		name    string
		entry   *WorkoutEntry
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
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       1,
				Reps:       1,
				LoadKg:     100.0,
				RPE:        8,
			},
			wantErr: false,
		},
		{
			name: "valid entry with all fields",
			entry: &WorkoutEntry{
				ID:                   uuid.New(),
				WorkoutID:            workoutID,
				Order:                0,
				ExerciseID:           exerciseID,
				DisplayName:          stringPtr("Bench Press"),
				EntryType:            stringPtr("main"),
				Sets:                 5,
				Reps:                 5,
				LoadKg:               100.0,
				RPE:                  7,
				EntryStart:           timePtr(now),
				EntryEnd:             timePtr(now.Add(30 * time.Minute)),
				PlannedRestSeconds:   intPtr(180),
				PerformedRestSeconds: intPtr(200),
				PerSetRestOverrides:  []int{90, 120, 120},
				Notes:                stringPtr("Felt strong"),
			},
			wantErr: false,
		},
		{
			name: "nil entry_type (valid - nullable)",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  nil,
				Sets:       1,
				Reps:       1,
				LoadKg:     100.0,
				RPE:        8,
			},
			wantErr: false,
		},
		{
			name: "custom entry_type (valid - user-defined)",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("custom_type"),
				Sets:       1,
				Reps:       1,
				LoadKg:     100.0,
				RPE:        8,
			},
			wantErr: false,
		},
		{
			name: "entry_type too long",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr(stringWithLength(51)),
				Sets:       1,
				Reps:       1,
				LoadKg:     100.0,
				RPE:        8,
			},
			wantErr: true,
			errCode: ErrCodeInvalidEntryType,
		},
		{
			name: "sets zero",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       0,
				Reps:       1,
				LoadKg:     100.0,
				RPE:        8,
			},
			wantErr: true,
		},
		{
			name: "sets negative",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       -1,
				Reps:       1,
				LoadKg:     100.0,
				RPE:        8,
			},
			wantErr: true,
		},
		{
			name: "reps zero",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       1,
				Reps:       0,
				LoadKg:     100.0,
				RPE:        8,
			},
			wantErr: true,
		},
		{
			name: "load_kg negative",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       1,
				Reps:       1,
				LoadKg:     -1.0,
				RPE:        8,
			},
			wantErr: true,
		},
		{
			name: "load_kg zero (bodyweight allowed)",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       1,
				Reps:       1,
				LoadKg:     0.0,
				RPE:        8,
			},
			wantErr: false,
		},
		{
			name: "RPE too low",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       1,
				Reps:       1,
				LoadKg:     100.0,
				RPE:        0,
			},
			wantErr: true,
			errCode: ErrCodeInvalidRPE,
		},
		{
			name: "RPE too high",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       1,
				Reps:       1,
				LoadKg:     100.0,
				RPE:        11,
			},
			wantErr: true,
			errCode: ErrCodeInvalidRPE,
		},
		{
			name: "order negative",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      -1,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       1,
				Reps:       1,
				LoadKg:     100.0,
				RPE:        8,
			},
			wantErr: true,
		},
		{
			name: "display_name too long",
			entry: &WorkoutEntry{
				ID:          uuid.New(),
				WorkoutID:   workoutID,
				Order:       0,
				ExerciseID:  exerciseID,
				DisplayName: stringPtr(stringWithLength(201)),
				EntryType:   stringPtr("top"),
				Sets:        1,
				Reps:        1,
				LoadKg:      100.0,
				RPE:         8,
			},
			wantErr: true,
		},
		{
			name: "video_object_key too long",
			entry: &WorkoutEntry{
				ID:             uuid.New(),
				WorkoutID:      workoutID,
				Order:          0,
				ExerciseID:     exerciseID,
				EntryType:      stringPtr("top"),
				Sets:           1,
				Reps:           1,
				LoadKg:         100.0,
				RPE:            8,
				VideoObjectKey: stringPtr(stringWithLength(501)),
			},
			wantErr: true,
		},
		{
			name: "load_kg rounding to 0.1kg",
			entry: &WorkoutEntry{
				ID:         uuid.New(),
				WorkoutID:  workoutID,
				Order:      0,
				ExerciseID: exerciseID,
				EntryType:  stringPtr("top"),
				Sets:       1,
				Reps:       1,
				LoadKg:     100.123, // Should round to 100.1
				RPE:        8,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateWorkoutEntry(tt.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWorkoutEntry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				if domainErr, ok := err.(*DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("ValidateWorkoutEntry() error code = %v, want %v", domainErr.Code, tt.errCode)
					}
				}
			}
			// Verify load rounding
			if !tt.wantErr && tt.entry != nil && tt.entry.LoadKg > 0 {
				rounded := RoundLoad(tt.entry.LoadKg)
				if tt.entry.LoadKg != rounded {
					t.Errorf("ValidateWorkoutEntry() should round load_kg, got %v, want %v", tt.entry.LoadKg, rounded)
				}
			}
		})
	}
}

func TestValidateWorkout(t *testing.T) {
	now := time.Now()
	past := now.Add(-1 * time.Hour)
	future := now.Add(1 * time.Hour)
	exerciseID := uuid.New()

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
		errCode string
	}{
		{
			name:    "nil workout",
			workout: nil,
			wantErr: true,
		},
		{
			name: "valid workout with minimal fields",
			workout: &Workout{
				ID:        uuid.New(),
				Timestamp: past,
				CreatedAt: now,
				UpdatedAt: now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  stringPtr("top"),
						Sets:       1,
						Reps:       1,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid workout with all fields",
			workout: &Workout{
				ID:             uuid.New(),
				Timestamp:      past,
				SessionStart:   timePtr(past),
				SessionEnd:     timePtr(past.Add(60 * time.Minute)),
				BodyWeightKg:   float64Ptr(74.5),
				FatigueLevel:   intPtr(2),
				SleepHours:     float64Ptr(7.5),
				ConditionNotes: stringPtr("Feeling good"),
				ProgramContext: []string{"Cycle 1", "Week 3", "Day 2"},
				Notes:          stringPtr("Great session"),
				CreatedAt:      now,
				UpdatedAt:      now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  stringPtr("top"),
						Sets:       1,
						Reps:       1,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "timestamp in future",
			workout: &Workout{
				ID:        uuid.New(),
				Timestamp: future,
				CreatedAt: now,
				UpdatedAt: now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  stringPtr("top"),
						Sets:       1,
						Reps:       1,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: true,
			errCode: ErrCodeInvalidTimestamp,
		},
		{
			name: "no entries",
			workout: &Workout{
				ID:        uuid.New(),
				Timestamp: past,
				CreatedAt: now,
				UpdatedAt: now,
				Entries:   []WorkoutEntry{},
			},
			wantErr: true,
		},
		{
			name: "session_start after session_end",
			workout: &Workout{
				ID:           uuid.New(),
				Timestamp:    past,
				SessionStart: timePtr(past.Add(60 * time.Minute)),
				SessionEnd:   timePtr(past),
				CreatedAt:    now,
				UpdatedAt:    now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  stringPtr("top"),
						Sets:       1,
						Reps:       1,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fatigue_level too low",
			workout: &Workout{
				ID:           uuid.New(),
				Timestamp:    past,
				FatigueLevel: intPtr(0),
				CreatedAt:    now,
				UpdatedAt:    now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  stringPtr("top"),
						Sets:       1,
						Reps:       1,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: true,
			errCode: ErrCodeInvalidFatigueLevel,
		},
		{
			name: "fatigue_level too high",
			workout: &Workout{
				ID:           uuid.New(),
				Timestamp:    past,
				FatigueLevel: intPtr(6),
				CreatedAt:    now,
				UpdatedAt:    now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  stringPtr("top"),
						Sets:       1,
						Reps:       1,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: true,
			errCode: ErrCodeInvalidFatigueLevel,
		},
		{
			name: "body_weight_kg zero or negative",
			workout: &Workout{
				ID:           uuid.New(),
				Timestamp:    past,
				BodyWeightKg: float64Ptr(0),
				CreatedAt:    now,
				UpdatedAt:    now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  stringPtr("top"),
						Sets:       1,
						Reps:       1,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sleep_hours negative",
			workout: &Workout{
				ID:         uuid.New(),
				Timestamp:  past,
				SleepHours: float64Ptr(-1),
				CreatedAt:  now,
				UpdatedAt:  now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  stringPtr("top"),
						Sets:       1,
						Reps:       1,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sleep_hours over 24",
			workout: &Workout{
				ID:         uuid.New(),
				Timestamp:  past,
				SleepHours: float64Ptr(25),
				CreatedAt:  now,
				UpdatedAt:  now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  stringPtr("top"),
						Sets:       1,
						Reps:       1,
						LoadKg:     100.0,
						RPE:        8,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "entry with nil entry_type (valid)",
			workout: &Workout{
				ID:        uuid.New(),
				Timestamp: past,
				CreatedAt: now,
				UpdatedAt: now,
				Entries: []WorkoutEntry{
					{
						ID:         uuid.New(),
						WorkoutID:  uuid.New(),
						Order:      0,
						ExerciseID: exerciseID,
						EntryType:  nil,
						Sets:       1,
						Reps:       1,
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
			err := ValidateWorkout(tt.workout)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateWorkout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errCode != "" {
				if domainErr, ok := err.(*DomainError); ok {
					if domainErr.Code != tt.errCode {
						t.Errorf("ValidateWorkout() error code = %v, want %v", domainErr.Code, tt.errCode)
					}
				}
			}
		})
	}
}

// Helper functions
func timePtr(t time.Time) *time.Time {
	return &t
}

func intPtr(i int) *int {
	return &i
}
