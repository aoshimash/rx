package domain

import (
	"github.com/google/uuid"
	"time"
)

// PlanSnapshot represents a snapshot of planned values at execution time.
type PlanSnapshot struct {
	ProgramNodeID      *uuid.UUID `json:"program_node_id,omitempty"`
	TargetSets         *int       `json:"target_sets,omitempty"`
	TargetReps         *int       `json:"target_reps,omitempty"`
	TargetRPE          *int       `json:"target_rpe,omitempty"`
	TargetLoadKg       *float64   `json:"target_load_kg,omitempty"`
	Percent1RM         *float64   `json:"percent_1rm,omitempty"`
	PlannedRestSeconds *int       `json:"planned_rest_seconds,omitempty"`
}

// WorkoutEntry represents a single performed exercise entry within a workout session.
type WorkoutEntry struct {
	ID                   uuid.UUID     `json:"id"`
	WorkoutID            uuid.UUID     `json:"workout_id"`
	Order                int           `json:"order"`
	ExerciseID           uuid.UUID     `json:"exercise_id"`
	DisplayName          *string       `json:"display_name,omitempty"`
	EntryType            string        `json:"entry_type"`
	Sets                 int           `json:"sets"`
	Reps                 int           `json:"reps"`
	LoadKg               float64       `json:"load_kg"`
	RPE                  int           `json:"rpe"`
	EntryStart           *time.Time    `json:"entry_start,omitempty"`
	EntryEnd             *time.Time    `json:"entry_end,omitempty"`
	PlannedRestSeconds   *int          `json:"planned_rest_seconds,omitempty"`
	PerformedRestSeconds *int          `json:"performed_rest_seconds,omitempty"`
	PerSetRestOverrides  []int         `json:"per_set_rest_overrides,omitempty"`
	ProgramNodeID        *uuid.UUID    `json:"program_node_id,omitempty"`
	PlanSnapshot         *PlanSnapshot `json:"plan_snapshot,omitempty"`
	Notes                *string       `json:"notes,omitempty"`
	VideoObjectKey       *string       `json:"video_object_key,omitempty"`
}

// Workout represents a completed training session containing performed entries.
type Workout struct {
	ID             uuid.UUID      `json:"id"`
	Timestamp      time.Time      `json:"timestamp"`
	SessionStart   *time.Time     `json:"session_start,omitempty"`
	SessionEnd     *time.Time     `json:"session_end,omitempty"`
	BodyWeightKg   *float64       `json:"body_weight_kg,omitempty"`
	FatigueLevel   *int           `json:"fatigue_level,omitempty"`
	SleepHours     *float64       `json:"sleep_hours,omitempty"`
	ConditionNotes *string        `json:"condition_notes,omitempty"`
	ProgramNodeID  *uuid.UUID     `json:"program_node_id,omitempty"`
	ProgramContext []string       `json:"program_context,omitempty"`
	Notes          *string        `json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Entries        []WorkoutEntry `json:"entries"`
}
