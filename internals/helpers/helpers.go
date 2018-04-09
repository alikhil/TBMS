package helpers

// Error is exception like simple class
type Error struct {
	Message string
}

func (e *Error) Error() string {
	return e.Message
}
