package errs

type AppError struct {
	Codes   int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}
func (e AppError) Code() int {
	return e.Codes
}

func NewNotFoundError(message string) error {
	return AppError{
		Message: message,
	}
}

func NewNotImplement(message string) error {
	return AppError{
		Message: message,
	}
}
