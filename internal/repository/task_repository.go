package repository

import (
	"context"
	"time"

	"github.com/siluk00/task_scheduler/internal/domain"
)

// contexts are used to pass deadlines, cancellation signals, and other request-scoped values across
// API boundaries and between processes.
// It's an interface with four methods:
// - Deadline: returns the time when the context will be canceled, if any.
// - Done: returns a channel that is closed when the context is canceled or times out.
// - Err: returns an error if the context has been canceled or timed out.
// - Value: retrieves a value associated with the context using a key.
type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	FindById(ctx context.Context, id string) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	List(ctx context.Context, status domain.TaskStatus) ([]*domain.Task, error)
	Delete(ctx context.Context, id string) error
	FindScheduled(ctx context.Context, from, to time.Time) ([]*domain.Task, error)
}
