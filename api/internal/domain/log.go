package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// LogEntry represents a single performed exercise entry within a training log.
type LogEntry struct {
	ID             uuid.UUID       `json:"id"`
	LogID          uuid.UUID       `json:"log_id"`
	Order          int             `json:"order"`
	ExerciseName   string          `json:"exercise_name"`
	Sets           *int            `json:"sets,omitempty"`
	Reps           *int            `json:"reps,omitempty"`
	LoadKg         *float64        `json:"load_kg,omitempty"`
	RPE            *int            `json:"rpe,omitempty"`
	Notes          *string         `json:"notes,omitempty"`
	VideoObjectKey *string         `json:"video_object_key,omitempty"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
}

// Log represents a record of actual training performed.
type Log struct {
	ID          uuid.UUID       `json:"id"`
	PlanID      *uuid.UUID      `json:"plan_id,omitempty"`
	PerformedAt time.Time       `json:"performed_at"`
	Notes       *string         `json:"notes,omitempty"`
	Metadata    json.RawMessage `json:"metadata,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Entries     []LogEntry      `json:"entries"`
}
