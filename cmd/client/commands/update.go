package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/siluk00/task_scheduler/internal/domain"
	"github.com/spf13/cobra"
)

func NewUpdateCommand() *cobra.Command {
	var (
		name        string
		description string
		command     string
		status      string
		scheduledAt string
		file        string
	)

	cmd := &cobra.Command{
		Use:   "update <task-id>",
		Short: "Update a task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskID := args[0]
			var task domain.Task
			var err error

			if file != "" {
				task, err = loadTaskFromFile(file)
				if err != nil {
					fmt.Printf("Error loading task from file: %v\n", err)
					return
				} else {
					resp, err := apiClient.Get(baseUrl + "/tasks/" + taskID)
					if err != nil {
						fmt.Printf("Error getting the task: %v\n", err)
						return
					}
					defer resp.Body.Close()

					if resp.StatusCode != http.StatusOK {
						body, _ := io.ReadAll(resp.Body)
						fmt.Printf("Error getting task: %s\n", string(body))
						return
					}

					if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
						fmt.Printf("Error decoding the task: %v\n", err)
						return
					}

					if name != "" {
						task.Name = name
					}
					if description != "" {
						task.Description = description
					}
					if command != "" {
						task.Command = command
					}
					if status != "" {
						task.Status = domain.TaskStatus(status)
					}
					if scheduledAt != "" {
						task.ScheduledAt, err = time.Parse(time.RFC3339, scheduledAt)
						if err != nil {
							fmt.Printf("Invalid scheduled_at format: %v\n", err)
							return
						}
					}
				}

				if err := updateTask(taskID, task); err != nil {
					fmt.Printf("Error updating task: %v\n", err)
					return
				}

				fmt.Println("Task updated successfully")
			}
		},
	}

	// Defining the string flags with --name -n
	cmd.Flags().StringVarP(&name, "name", "n", "", "Task name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Task description")
	cmd.Flags().StringVarP(&command, "command", "c", "", "Command to execute")
	cmd.Flags().StringVarP(&status, "status", "s", "", "Task status")
	cmd.Flags().StringVarP(&scheduledAt, "scheduled-at", "t", "", "Scheduled time (RFC3339 format)")
	cmd.Flags().StringVarP(&file, "file", "f", "", "JSON file with task data")

	return cmd
}

func updateTask(taskID string, task domain.Task) error {
	task.ID = taskID
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("error marshalling task: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, baseUrl+"/tasks/"+taskID, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-type", "application/json")

	resp, err := apiClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body ))
	}

	return nil
}
