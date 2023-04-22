package errors

import "fmt"

var (
	ErrStreamAlreadyExists  = fmt.Errorf("stream already exists")
	ErrInvalidStreamVersion = fmt.Errorf("invalid stream version")
)
