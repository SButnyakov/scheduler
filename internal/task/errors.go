package task

import "errors"

var (
	ErrInvalidState    = errors.New("failed to set invalid state")
	ErrInvalidType     = errors.New("failed to set invalid type")
	ErrInvalidPriority = errors.New("failed to set invalid priority")
)
