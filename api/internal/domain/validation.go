package domain

import (
	"fmt"
	"math"
	"time"
)

// RoundLoad rounds a weight value to 0.1kg precision.
func RoundLoad(kg float64) float64 {
	return math.Round(kg*10) / 10
}

// ValidateRPE checks if RPE is in the valid range (1-10).
func ValidateRPE(rpe int) error {
	if rpe < 1 || rpe > 10 {
		return &DomainError{
			Code:    ErrCodeInvalidRPE,
			Message: "RPE must be between 1 and 10",
			Details: map[string]interface{}{
				"value": rpe,
				"min":   1,
				"max":   10,
			},
		}
	}
	return nil
}

// ValidateFatigueLevel checks if fatigue level is in the valid range (1-5).
func ValidateFatigueLevel(level int) error {
	if level < 1 || level > 5 {
		return &DomainError{
			Code:    ErrCodeInvalidFatigueLevel,
			Message: "Fatigue level must be between 1 and 5",
			Details: map[string]interface{}{
				"value": level,
				"min":   1,
				"max":   5,
			},
		}
	}
	return nil
}

// ValidateTimestamp checks if timestamp is not in the future.
func ValidateTimestamp(t time.Time) error {
	now := time.Now()
	if t.After(now) {
		return &DomainError{
			Code:    ErrCodeInvalidTimestamp,
			Message: "Timestamp cannot be in the future",
			Details: map[string]interface{}{
				"timestamp": t,
				"now":       now,
			},
		}
	}
	return nil
}

// ValidateEntryType checks if entry_type is one of the valid values.
func ValidateEntryType(entryType string) error {
	validTypes := map[string]bool{
		"top":       true,
		"main":      true,
		"backoff":   true,
		"accessory": true,
	}
	if !validTypes[entryType] {
		return &DomainError{
			Code:    ErrCodeInvalidEntryType,
			Message: "Entry type must be one of: top, main, backoff, accessory",
			Details: map[string]interface{}{
				"value":      entryType,
				"validTypes": []string{"top", "main", "backoff", "accessory"},
			},
		}
	}
	return nil
}

// ValidateRequiredString checks if a required string field is not empty.
func ValidateRequiredString(field, value string) error {
	if value == "" {
		return &ValidationError{
			Field:   field,
			Message: "required field cannot be empty",
		}
	}
	return nil
}

// ValidateStringLength checks if a string is within the specified length range.
func ValidateStringLength(field, value string, min, max int) error {
	length := len(value)
	if length < min || length > max {
		return &ValidationError{
			Field:   field,
			Message: fmt.Sprintf("string length must be between %d and %d characters", min, max),
		}
	}
	return nil
}
