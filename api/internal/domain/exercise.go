package domain

import (
	"github.com/google/uuid"
	"time"
)

// Exercise represents a catalog entry for a canonical exercise.
type Exercise struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	Aliases       []string  `json:"aliases,omitempty"`
	MuscleGroups  []string  `json:"muscle_groups,omitempty"`
	LoadIncrement *float64  `json:"load_increment,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
