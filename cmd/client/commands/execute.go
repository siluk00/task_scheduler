package commands

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

func NewExecuteCommand() *cobra.Command {
	var (
		async bool
	)
	cmd := &cobra.Command{
		Use:   "execute <task-id>",
		Short: "Execute a task manually",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			taskID := args[0]
			url := baseUrl + "/tasks/" + taskID + "/execute"

			if async {
				url += "?async=true"
			}

			req, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				fmt.Printf("Error creating task: %v\n,", err)
				return
			}

			resp, err := apiClient.Do(req)
			if err != nil {
				fmt.Printf("Error executing task: %v\n", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error executing task: status %d\n", resp.StatusCode)
				return
			}

			if async {
				fmt.Println("Task execution started asynchronously")
			} else {
				fmt.Println("task executed succesfully")
			}
		},
	}

	cmd.Flags().BoolVarP(&async, "async", "a", false, "Execute asynchronously")

	return cmd
}
