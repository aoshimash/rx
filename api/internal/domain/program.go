package domain

import (
	"time"

	"github.com/google/uuid"
)

// FieldDef defines a field schema for program or log entries.
type FieldDef struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Options     []string `json:"options,omitempty"`
	Description string   `json:"description,omitempty"`
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
	ID           uuid.UUID             `json:"id"`
	ProgramID    uuid.UUID             `json:"program_id"`
	GroupID      *uuid.UUID            `json:"group_id,omitempty"`
	FieldGroupID *uuid.UUID            `json:"field_group_id,omitempty"`
	SessionName  string                `json:"session_name"`
	Order        int                   `json:"order"`
	Date         *DateOnly             `json:"date,omitempty"`
	Entries      []ProgramSessionEntry `json:"entries,omitempty"`
}

// Program represents a reusable training program template with embedded sessions.
type Program struct {
	ID        uuid.UUID        `json:"id"`
	Name      string           `json:"name"`
	Notes     *string          `json:"notes,omitempty"`
	Groups    []ProgramGroup   `json:"groups,omitempty"`
	Sessions  []ProgramSession `json:"sessions,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}
