package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	var (
		status       string
		outputFormat string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		Run: func(cmd *cobra.Command, args []string) {
			url := baseUrl + "/tasks"
			if status != "" {
				url += "?status=" + status
			}

			resp, err := apiClient.Get(url)
			if err != nil {
				fmt.Printf("Error making request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				body, _ := io.ReadAll(resp.Body)
				fmt.Printf("Error listing tasks: %s\n", string(body))
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response: %v\n", err)
				return
			}

			if outputFormat == "json" {
				fmt.Println(string(body))
			} else {
				var tasks []map[string]interface{}
				if err = json.Unmarshal(body, &tasks); err != nil {
					fmt.Printf("Error decoding response: %v\n", err)
					return
				}
				printTasksPretty(tasks)
			}
		},
	}

	cmd.Flags().StringVarP(&status, "status", "s", "", "Filter by status (pending, running, completed, failed)")
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "pretty", "format (pretty|json)")

	return cmd
}

func printTasksPretty(tasks []map[string]interface{}) {
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return
	}

	fmt.Printf("Found %d tasks:\n", len(tasks))
	for i, task := range tasks {
		fmt.Printf("\nTask #%d:\n", i+1)
		fmt.Printf("\tId:%s\n", task["id"])
		fmt.Printf("\tName: %s\n", task["name"])
		fmt.Printf("\tStatus: %s\n", task["status"])
		if task["scheduled_at"] != nil {
			fmt.Printf("\tScheduled At: %s\n", task["scheduled_at"])
		}
	}
}
