package utils

type AppError struct {
	Code       string //client facing error status
	Message    string
	Details    string //additional details about the error
	StatusCode int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Details + ": " + e.Err.Error()
	}
	return e.Details
}

func NewAppError(statusCode int, code, message string, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Err:        err,
	}
}

var (
	ErrBadRequest         = NewAppError(400, "BAD_REQUEST", "Invalid request body", nil)
	ErrUnauthorized       = NewAppError(401, "UNAUTHORIZED", "Request unauthorized", nil)
	ErrInvalidCredentials = NewAppError(401, "INVALID_CREDENTIALS", "Invalid credentials", nil)
	ErrTokenGenFailure    = NewAppError(500, "TOKEN_GEN_FAILED", "Failed to generate authentication token", nil)
	ErrTokenInvalid       = NewAppError(401, "INVALID_TOKEN", "Invalid authentication token", nil)
	ErrTokenMissing       = NewAppError(400, "TOKEN_MISSING", "Authentication token missing", nil)
	ErrInternalServer     = NewAppError(500, "INTERNAL_ERROR", "Internal server error", nil)
	ErrNotFound           = NewAppError(500, "ENTITY_NOT_FOUND", "Requested entity not found", nil)
)

func (e *AppError) WithDetails(details string) *AppError {
	newErr := *e
	newErr.Details = details
	return &newErr
}