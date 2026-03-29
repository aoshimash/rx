package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// LogEntry represents a single performed exercise entry within a training log.
type LogEntry struct {
	ID             uuid.UUID              `json:"id"`
	LogID          uuid.UUID              `json:"log_id"`
	Order          int                    `json:"order"`
	ExerciseName   string                 `json:"exercise_name"`
	Fields         map[string]interface{} `json:"fields,omitempty"`
	Notes          *string                `json:"notes,omitempty"`
	VideoObjectKey *string                `json:"video_object_key,omitempty"`
	StartedAt      *time.Time             `json:"started_at,omitempty"`
	FinishedAt     *time.Time             `json:"finished_at,omitempty"`
}

// LoggedSession represents a session that has a recorded log, used for navigation.
type LoggedSession struct {
	SessionName string     `json:"session_name"`
	LogID       uuid.UUID  `json:"log_id"`
	PerformedAt time.Time  `json:"performed_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
}

// Log represents a record of actual training performed.
type Log struct {
	ID          uuid.UUID       `json:"id"`
	ProgramID   *uuid.UUID      `json:"program_id,omitempty"`
	SessionName *string         `json:"session_name,omitempty"`
	PerformedAt time.Time       `json:"performed_at"`
	StartedAt   *time.Time      `json:"started_at,omitempty"`
	FinishedAt  *time.Time      `json:"finished_at,omitempty"`
	Notes       *string         `json:"notes,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Entries     []LogEntry      `json:"entries"`
}
