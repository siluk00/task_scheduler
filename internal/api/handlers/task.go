// @title Task Scheduler API
// @version 1.0
// @deescription API for task scheduling
// @host localhost:8080
// @BasePath /

package handlers

import (
	"github.com/siluk00/task_scheduler/internal/repository"
)

type taskHandler struct {
	repo repository.TaskHandler
}

func NewTaskHandler(repo repository.TaskHandler) *taskHandler {
	return &taskHandler{
		repo: repo,
	}
}
