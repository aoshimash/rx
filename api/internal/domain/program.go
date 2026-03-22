package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// ProgramStatus represents the lifecycle state of a Program.
type ProgramStatus string

const (
	ProgramStatusActive    ProgramStatus = "active"
	ProgramStatusCompleted ProgramStatus = "completed"
	ProgramStatusPlanned   ProgramStatus = "planned"
)

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
// Status transitions from "active" to "completed" when all sessions have been logged.
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
