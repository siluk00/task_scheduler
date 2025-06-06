package commands

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

func NewGetCommand() *cobra.Command {
	var (
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "get <task-id>",
		Short: "Get task details",
		Args:  cobra.ExactArgs(1), // Ensure exactly one argument is provided
		Run: func(cmd *cobra.Command, args []string) {
			taskID := args[0]

			resp, err := apiClient.Get(baseUrl + "/tasks/" + taskID)
			if err != nil {
				cmd.Printf("Error making request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Error getting task: %s\n", string(body))
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				cmd.Printf("Error reading response body: %v\n", err)
				return
			}

			if outputFormat == "json" {
				fmt.Println(string(body)) // Print JSON response directly
			} else {
				var task map[string]interface{}
				if err := json.Unmarshal(body, &task); err != nil {
					cmd.Printf("Error parsing JSON response: %v\n", err)
					return
				}
				printTaskPretty(task)
			}
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "pretty", "Output format(json|pretty)")
	return cmd
}

func printTaskPretty(task map[string]interface{}) {
	fmt.Println("Task Details:")
	fmt.Printf("ID:\t\t %s\n", task["id"])
	fmt.Printf("Name:\t\t %s\n", task["name"])
	fmt.Printf("Description:\t %s\n", task["description"])
	fmt.Printf("Command:\t %s\n", task["command"])
	fmt.Printf("Status:\t\t %s\n", task["status"])
	fmt.Printf("Created At:\t %s\n", task["created_at"])
	fmt.Printf("Updated At:\t %s\n", task["updated_at"])
	if scheduledAt, ok := task["scheduled_at"]; ok {
		fmt.Printf("Scheduled At:\t %s\n", scheduledAt)
	} else {
		fmt.Println("Scheduled At:\t Not Scheduled")
	}
}
