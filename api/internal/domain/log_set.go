package domain

import "github.com/google/uuid"

// LogSet represents a single set performed within a log entry.
type LogSet struct {
	ID             uuid.UUID              `json:"id"`
	EntryID        uuid.UUID              `json:"entry_id"`
	SetNumber      int                    `json:"set_number"`
	Fields         map[string]interface{} `json:"fields"`
	VideoObjectKey *string                `json:"video_object_key,omitempty"`
	Notes          *string                `json:"notes,omitempty"`
}
