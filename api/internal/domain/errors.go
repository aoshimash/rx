package domain

// ValidationError represents a field-level validation error.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

// DomainError represents a domain-level error with error code and details.
type DomainError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

func (e *DomainError) Error() string {
	return e.Code + ": " + e.Message
}

// Error codes
const (
	ErrCodeInvalidTimestamp     = "INVALID_TIMESTAMP"
	ErrCodeInvalidRPE           = "INVALID_RPE"
	ErrCodeInvalidFatigueLevel  = "INVALID_FATIGUE_LEVEL"
	ErrCodeMissingRequiredField = "MISSING_REQUIRED_FIELD"
	ErrCodeInvalidEntryType     = "INVALID_ENTRY_TYPE"
)
