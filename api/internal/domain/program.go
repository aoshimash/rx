package domain

import (
	"github.com/google/uuid"
	"time"
)

// ProgramNode represents a node in the program tree (cycle, week, day, block, or exercise prescription).
type ProgramNode struct {
	ID                uuid.UUID     `json:"id"`
	ProgramID         uuid.UUID     `json:"program_id"`
	ParentID          *uuid.UUID    `json:"parent_id,omitempty"`
	Name              string        `json:"name"`
	NodeType          string        `json:"node_type"`
	Order             int           `json:"order"`
	Children          []ProgramNode `json:"children,omitempty"`
	// Prescription fields (for leaf nodes)
	ExerciseID         *uuid.UUID `json:"exercise_id,omitempty"`
	TargetSets         *int       `json:"target_sets,omitempty"`
	TargetReps         *int       `json:"target_reps,omitempty"`
	TargetRPE          *int       `json:"target_rpe,omitempty"`
	Percent1RM         *float64   `json:"percent_1rm,omitempty"`
	PlannedRestSeconds *int       `json:"planned_rest_seconds,omitempty"`
	MuscleGroups       []string   `json:"muscle_groups,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
}

// Program represents a training program containing a recursive tree of nodes.
type Program struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Description *string       `json:"description,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	RootNodes   []ProgramNode `json:"root_nodes,omitempty"`
}
