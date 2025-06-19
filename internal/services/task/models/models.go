package models

import (
	"fmt"
	"time"
)

type TaskStatus string

const (
	TaskPending TaskStatus = "pending"
	TaskRunning TaskStatus = "running"
	TaskDone    TaskStatus = "done"
	TaskFailed  TaskStatus = "failed"
)

var ValidStatuses = []TaskStatus{
	TaskPending,
	TaskRunning,
	TaskDone,
	TaskFailed,
}

type Task struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt time.Time  `json:"completed_at"`
}

func IsValidStatus(status TaskStatus) bool {
	for _, validStatus := range ValidStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}
func (t *Task) SetStatus(status TaskStatus) error {
	if !IsValidStatus(status) {
		return fmt.Errorf("invalid task status: %s", status)
	}
	t.Status = status
	return nil
}
