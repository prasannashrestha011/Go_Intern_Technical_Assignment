package utils

//custom error type
// status->true(success)
// status->false(failed)
type AppError struct {
	Message string
	Err error
	Code int
}

func (e *AppError) Error() string{
	return e.Message
}

func New(code int,message string,err error) *AppError{

	return &AppError{
		Message: message,
		Code: code,
		Err:err,
		}
}

var (
	ErrBadRequest=New(400,"Invalid request body",nil)
	ErrUnAuthorized=New(401,"Request unauthorized",nil)
	ErrInvalidCredentials=New(401,"Invalid credentials",nil)
	ErrTokenGenFailure=New(500,"Failed to generate authentication token",nil)
	ErrTokenInvalid=New(401,"Invalid authentication token",nil)
	ErrTokenMissing=New(400,"Authentication token missing!!",nil)
	ErrInternalServerError=New(500,"Internal server error",nil)
)
