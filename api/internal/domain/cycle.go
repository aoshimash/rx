package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Cycle represents a single conversion of a Program into Plans.
// It groups Plans generated from one conversion, enabling clear traceability.
type Cycle struct {
	ID        uuid.UUID       `json:"id"`
	ProgramID uuid.UUID       `json:"program_id"`
	Name      string          `json:"name"`
	Notes     *string         `json:"notes,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	CreatedAt time.Time       `json:"created_at"`
}

// ValidateCycle validates a Cycle entity.
func ValidateCycle(c *Cycle) error {
	if c == nil {
		return &ValidationError{
			Field:   "cycle",
			Message: "cycle cannot be nil",
		}
	}

	// Validate required fields
	if err := ValidateRequiredString("name", c.Name); err != nil {
		return err
	}
	if err := ValidateStringLength("name", c.Name, 1, 200); err != nil {
		return err
	}

	// Validate notes length if provided
	if c.Notes != nil {
		if err := ValidateStringLength("notes", *c.Notes, 0, 5000); err != nil {
			return err
		}
	}

	return nil
}
