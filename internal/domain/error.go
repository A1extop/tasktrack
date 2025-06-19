package domain

import "errors"

var (
	ErrCannotDeleteTask  = errors.New("cannot delete task: task is running")
	ErrNamesTaskIsEmpty  = errors.New("task name is empty")
	ErrNameTooLong       = errors.New("task name is too long (max 100 characters)")
	ErrTaskAlreadyExists = errors.New("task with this name already exists")
	ErrInvalidID         = errors.New("invalid task ID")
	ErrDataIsEmpty       = errors.New("no data provided to update task")
	ErrInvalidStatus     = errors.New("invalid task status")
	ErrImmutableField    = errors.New("createdAt field cannot be changed")
	ErrTaskNotFound      = errors.New("task not found")
)
