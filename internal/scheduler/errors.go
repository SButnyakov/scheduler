package scheduler

import "errors"

var (
	ErrNoTasksLeft     = errors.New("no tasks")
	ErrNoSpaceForReady = errors.New("no space for new ready task")
)
