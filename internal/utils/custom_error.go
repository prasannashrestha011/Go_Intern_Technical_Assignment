package utils

//custom error type
// status->true(success)
// status->false(failed)
type CustomError struct {
	Status  bool
	Message string
}

func (e *CustomError) Error() string{
	return e.Message
}