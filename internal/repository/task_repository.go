package repository

import (
	"context"
	"time"

	"github.com/siluk00/task_scheduler/internal/domain"
)

// The interface for CRUD of the tasks
type TaskHandler interface {
	Create(ctx context.Context, task *domain.Task) error
	FindById(ctx context.Context, id string) (*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	List(ctx context.Context, status domain.TaskStatus) ([]*domain.Task, error)
	Delete(ctx context.Context, id string) error
	FindScheduled(ctx context.Context, from, to time.Time) ([]*domain.Task, error)
}
