package domain

import (
	"time"

	"github.com/google/uuid"
)

// FieldGroup defines a reusable set of field schemas for program and log entries.
// Each group contains paired program_fields and log_fields definitions.
type FieldGroup struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	ProgramFields []FieldDef `json:"program_fields"`
	LogFields     []FieldDef `json:"log_fields"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
