package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/siluk00/task_scheduler/internal/domain"
)

// Consumes eveything in the task_queue queue
func (w *TaskWorker) StartConsumer(ctx context.Context) error {
	msgs, err := w.msgQueue.Consume("task_queue")
	if err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}

	processor := NewTaskProcessor(w.taskRepo)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return errors.New("message channel closed")
			}

			var task domain.Task
			if err := json.Unmarshal(msg.Body, &task); err != nil {
				log.Printf("Failed to unmarshall task: %v", err)
				_ = msg.Nack(false, false)
				continue
			}

			log.Printf("Processing task %s", task.ID)
			if err := processor.ProcessTask(ctx, &task); err != nil {
				log.Printf("Failed to process task %s: %v", task.ID, err)
				_ = msg.Nack(false, true)
				continue
			}

			if err := msg.Ack(false); err != nil {
				log.Printf("failed to ack message %v", err)
			}
		}
	}
}
