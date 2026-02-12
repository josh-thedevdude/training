package helpers

type AppError struct {
	Code    int
	Message string
	Tag     string
}

func (e *AppError) Error() string {
	return e.Message
}
