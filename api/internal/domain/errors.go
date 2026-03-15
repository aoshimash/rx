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
	ErrCodeMissingRequiredField = "MISSING_REQUIRED_FIELD"
)

// Predefined domain errors
var (
	ErrNotFound = &DomainError{
		Code:    "NOT_FOUND",
		Message: "Resource not found",
	}
)

// Error codes for HTTP responses
const (
	ErrorCodeUnauthorized    = "UNAUTHORIZED"
	ErrorCodeValidationError = "VALIDATION_ERROR"
	ErrorCodeNotFound        = "NOT_FOUND"
	ErrorCodeConflict        = "CONFLICT"
	ErrorCodeInternalError   = "INTERNAL_ERROR"
)
