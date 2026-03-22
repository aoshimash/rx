package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// ProgramStatus represents the lifecycle state of a Program.
type ProgramStatus string

const (
	ProgramStatusCreated   ProgramStatus = "created"
	ProgramStatusOngoing   ProgramStatus = "ongoing"
	ProgramStatusCompleted ProgramStatus = "completed"
	ProgramStatusCancelled ProgramStatus = "cancelled"
)

// ValidateProgramStatusTransition checks if a status transition is allowed.
// Allowed: created→ongoing, ongoing→completed, ongoing→cancelled.
func ValidateProgramStatusTransition(from, to ProgramStatus) error {
	switch {
	case from == ProgramStatusCreated && to == ProgramStatusOngoing:
		return nil
	case from == ProgramStatusOngoing && to == ProgramStatusCompleted:
		return nil
	case from == ProgramStatusOngoing && to == ProgramStatusCancelled:
		return nil
	default:
		return &ValidationError{
			Field:   "status",
			Message: fmt.Sprintf("invalid status transition from '%s' to '%s'", from, to),
		}
	}
}

// ProgramSessionEntry represents a single exercise prescription within a program session.
type ProgramSessionEntry struct {
	ID           uuid.UUID       `json:"id"`
	SessionID    uuid.UUID       `json:"session_id"`
	Order        int             `json:"order"`
	ExerciseName string          `json:"exercise_name"`
	Sets         *int            `json:"sets,omitempty"`
	Reps         *int            `json:"reps,omitempty"`
	LoadKg       *float64        `json:"load_kg,omitempty"`
	RPE          *int            `json:"rpe,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	Metadata     json.RawMessage `json:"metadata,omitempty"`
}

// ProgramSession represents a single training session within a concrete Program.
type ProgramSession struct {
	ID          uuid.UUID             `json:"id"`
	ProgramID   uuid.UUID             `json:"program_id"`
	SessionName string                `json:"session_name"`
	Order       int                   `json:"order"`
	Date        *DateOnly             `json:"date,omitempty"`
	Entries     []ProgramSessionEntry `json:"entries,omitempty"`
}

// Program represents a concrete, immutable training program with embedded sessions.
// Generated from a ProgramTemplate or created manually.
// Status transitions: created → ongoing → completed/cancelled. All transitions are explicit user actions.
type Program struct {
	ID                uuid.UUID        `json:"id"`
	ProgramTemplateID *uuid.UUID       `json:"program_template_id,omitempty"`
	Name              string           `json:"name"`
	Status            ProgramStatus    `json:"status"`
	Notes             *string          `json:"notes,omitempty"`
	Metadata          json.RawMessage  `json:"metadata,omitempty"`
	Sessions          []ProgramSession `json:"sessions,omitempty"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}
