package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

func NewScheduleCommand() *cobra.Command {
	var (
		from string
		to   string
	)

	cmd := &cobra.Command{
		Use:   "scheduled",
		Short: "List scheduled tasks",
		Run: func(cmd *cobra.Command, args []string) {
			now := time.Now()
			defaultFrom := now.Format(time.RFC3339)
			defaultTo := now.Add(24 * time.Hour).Format(time.RFC3339)

			if from == "" {
				from = defaultFrom
			}
			if to == "" {
				to = defaultTo
			}

			url := fmt.Sprintf("%s/tasks/scheduled?from=%s&to=%s", baseUrl, from, to)
			resp, err := apiClient.Get(url)
			if err != nil {
				fmt.Printf("Error making request: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("Error reading response: %v\n", err)
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response: %v\n", err)
				return
			}

			var tasks []map[string]interface{}
			if err := json.Unmarshal(body, &tasks); err != nil {
				fmt.Printf("Error decoding response: %v\n", err)
				return
			}

			printScheduledTasks(tasks, from, to)
		},
	}

	cmd.Flags().StringVarP(&from, "from", "f", "", "Start time (RFC3339 format)")
	cmd.Flags().StringVarP(&to, "to", "t", "", "End time (RFC3339 format)")

	return cmd
}

func printScheduledTasks(tasks []map[string]interface{}, from, to string) {
	fmt.Printf("Scheduled tasks between %s and %s:\n", from, to)
	if len(tasks) == 0 {
		fmt.Println("No tasks found in this time range")
		return
	}

	for i, task := range tasks {
		fmt.Printf("\nTask #%d:\n", i+1)
		fmt.Printf("\tID: %s\n", task["id"])
		fmt.Printf("\tName: %s\n", task["name"])
		fmt.Printf("\tScheduled At:%s\n", task["scheduled_at"])
		fmt.Printf("\tStatus: %s\n", task["status"])
	}
}
