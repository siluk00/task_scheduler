package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/siluk00/task_scheduler/internal/domain"
	"github.com/spf13/cobra"
)

func NewCreateCommand() *cobra.Command {
	var (
		id          string
		name        string
		description string
		command     string
		status      string
		scheduledAt string
		file        string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new task",
		//Run is the function that will be executed when the command is called.
		//It receives the command and arguments as parameters.
		Run: func(cmd *cobra.Command, args []string) {
			var task domain.Task
			var err error

			if file != "" {
				task, err = loadTaskFromFile(file)
				if err != nil {
					fmt.Printf("Error loading task from file: %v\n", err)
					return
				}

			} else {
				task = domain.Task{
					ID:          id,
					Name:        name,
					Description: description,
					Command:     command,
					Status:      domain.TaskStatus(status),
				}

				if scheduledAt != "" {
					//time.RFC3339 is the format "2006-01-02T15:04:05Z07:00"
					//time.Parse parses a formatted string and returns the time.Time value.
					task.ScheduledAt, err = time.Parse(time.RFC3339, scheduledAt)
					if err != nil {
						fmt.Printf("Invalid scheduled-at format: %v\n", err)
						return
					}
				}
			}

			if err := createTask(task); err != nil {
				fmt.Printf("Error creating task: %v\n", err)
				return
			}

			fmt.Println("Task created successfully")
		},
	}

	//Each flag is a command-line option that can be passed to the command.
	//Flags are used to provide additional information or options to the command.
	//StringVarP binds a string flag to a variable.
	// The first argument is the variable to bind, the second is the flag name,
	// the third is the shorthand flag (optional), the fourth is the default value,
	// and the fifth is the usage description.
	// It has to be used in terminal like this:
	// taskctl create --id <task_id> --name <task_name> --command <task_command>

	cmd.Flags().StringVarP(&id, "id", "i", "", "Task ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Task name")
	cmd.Flags().StringVarP(&description, "description", "d", "", "Task description")
	cmd.Flags().StringVarP(&command, "command", "c", "", "Command to execute")
	cmd.Flags().StringVarP(&status, "status", "s", "pending", "Task status (pending, running, completed, failed)")
	cmd.Flags().StringVarP(&scheduledAt, "scheduled-at", "t", "", "Scheduled time in RFC3339 format")
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path to JSON file containing task data")

	//self-documented
	cmd.MarkFlagRequired("id")
	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("command")

	return cmd
}

// Posts a new task to the API server
func createTask(task domain.Task) error {
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	resp, err := apiClient.Post(baseUrl+"/tasks", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

func loadTaskFromFile(filePath string) (domain.Task, error) {
	var task domain.Task

	file, err := os.Open(filePath)
	if err != nil {
		return task, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&task); err != nil {
		return task, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return task, nil
}
