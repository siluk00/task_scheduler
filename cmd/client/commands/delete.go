package commands

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func NewDeleteCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "delete <task-id>",
		Short: "delete task",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskId := args[0]

			if err := deleteTask(taskId); err != nil {
				fmt.Printf("Error deleting task: %v\n", err)
				return
			}

			fmt.Println("Task deleted successfully")
		},
	}

	return cmd
}

func deleteTask(taskId string) error {
	req, err := http.NewRequest(http.MethodDelete, baseUrl+"/"+taskId, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	resp, err := apiClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}
