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
// Allowed: created→ongoing, ongoing→completed, ongoing→cancelled, cancelled→ongoing.
func ValidateProgramStatusTransition(from, to ProgramStatus) error {
	switch {
	case from == ProgramStatusCreated && to == ProgramStatusOngoing:
		return nil
	case from == ProgramStatusOngoing && to == ProgramStatusCompleted:
		return nil
	case from == ProgramStatusOngoing && to == ProgramStatusCancelled:
		return nil
	case from == ProgramStatusCancelled && to == ProgramStatusOngoing:
		return nil
	default:
		return &ValidationError{
			Field:   "status",
			Message: fmt.Sprintf("invalid status transition from '%s' to '%s'", from, to),
		}
	}
}

// FieldDef defines a field schema for program or log entries.
type FieldDef struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Options []string `json:"options,omitempty"`
}

// ProgramSessionEntry represents a single exercise prescription within a program session.
type ProgramSessionEntry struct {
	ID           uuid.UUID              `json:"id"`
	SessionID    uuid.UUID              `json:"session_id"`
	Order        int                    `json:"order"`
	ExerciseName string                 `json:"exercise_name"`
	Fields       map[string]interface{} `json:"fields,omitempty"`
	Notes        *string                `json:"notes,omitempty"`
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

// Program represents a concrete training program with embedded sessions.
// Status transitions: created → ongoing → completed/cancelled, cancelled → ongoing. All transitions are explicit user actions.
type Program struct {
	ID            uuid.UUID        `json:"id"`
	Name          string           `json:"name"`
	Status        ProgramStatus    `json:"status"`
	Notes         *string          `json:"notes,omitempty"`
	Metadata      json.RawMessage  `json:"metadata,omitempty"`
	ProgramFields []FieldDef       `json:"program_fields,omitempty"`
	LogFields     []FieldDef       `json:"log_fields,omitempty"`
	Sessions      []ProgramSession `json:"sessions,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}
