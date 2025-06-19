package worker

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/siluk00/task_scheduler/internal/domain"
	"github.com/siluk00/task_scheduler/internal/repository"
)

type TaskProcessor struct {
	taskRepo repository.TaskRepository
}

func NewTaskProcessor(repo repository.TaskRepository) *TaskProcessor {
	return &TaskProcessor{
		taskRepo: repo,
	}
}

func (p *TaskProcessor) ProcessTask(ctx context.Context, task *domain.Task) error {
	if task.Status != domain.TaskStatusRunning {
		task.Status = domain.TaskStatusRunning
		task.UpdatedAt = time.Now()
		if err := p.taskRepo.Update(ctx, task); err != nil {
			return fmt.Errorf("failed to update task: %v", err)
		}
	}

	output, err := p.executeCommand(task.Command)
	if err != nil {
		task.Status = domain.TaskStatusFailed
		log.Printf("Task %s failed: %v, Output: %s", task.ID, err, output)
	} else {
		task.Status = domain.TaskStatusCompleted
		log.Printf("Task %s completed succesfully, output: %s", task.ID, output)
	}

	task.UpdatedAt = time.Now()
	if err := p.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("failed to update task status: %v", err)
	}

	return nil
}

func (p *TaskProcessor) executeCommand(cmdString string) (string, error) {
	cmd := exec.Command("sh", "-c", cmdString)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
