package domain

import (
	"errors"
	"regexp"
	"time"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Command     string     `json:"command"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ScheduledAt time.Time  `json:"scheduled_at,omitempty"`
}

var (
	validStatuses = map[TaskStatus]bool{
		TaskStatusPending:   true,
		TaskStatusRunning:   true,
		TaskStatusCompleted: true,
		TaskStatusFailed:    true,
	}

	ErrInvalidTaskId      = errors.New("invalid task ID")
	ErrInvalidTaskStatus  = errors.New("invalid task status")
	ErrInvalidTaskName    = errors.New("invalid task name")
	ErrInvalidCommand     = errors.New("invalid command")
	ErrInvalidScheduledAt = errors.New("invalid scheduled time")
	//ErrTaskNotFound = errors.New("task not found")
	//ErrTaskAlreadyExists = errors.New("task already exists")
	//ErrTaskCreationFailed = errors.New("task creation failed")
	//ErrTaskUpdateFailed = errors.New("task update failed")
	//ErrTaskDeletionFailed = errors.New("task deletion failed")
	//ErrTaskSchedulingFailed = errors.New("task scheduling failed"
	//ErrTaskExecutionFailed = errors.New("task execution failed"
	//ErrTaskTimeout = errors.New("task execution timed out"
	//ErrTaskValidationFailed = errors.New("task validation failed"
	//ErrTaskCommandEmpty = errors.New("task command cannot be empty"
	//ErrTaskNameEmpty = errors.New("task name cannot be empty"
)

func (t *Task) Validate() error {
	if t.ID == "" || !isValidTaskId(t.ID) {
		return ErrInvalidTaskId
	}

	if t.Name == "" || len(t.Name) > 100 {
		return ErrInvalidTaskName
	}

	if t.Command == "" {
		return ErrInvalidCommand
	}

	if !validStatuses[t.Status] {
		return ErrInvalidTaskStatus
	}

	if !t.ScheduledAt.IsZero() && t.ScheduledAt.Before(time.Now()) {
		return ErrInvalidScheduledAt
	}

	return nil
}

func isValidTaskId(id string) bool {
	//O pacote regexp é usado para trabalhar com expressões regulares em Go.
	// A expressão regular `^[a-zA-Z0-9-]+$` verifica se o ID contém apenas letras,
	// números e hífens.
	match, _ := regexp.MatchString(`^[a-zA-Z0-9-]+$`, id)
	return match
}
