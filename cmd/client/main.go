package main

import (
	"log"
	"os"

	"github.com/siluk00/task_scheduler/cmd/client/commands"
	"github.com/spf13/cobra"
)

func main() {
	// A cobra.Command is the root command for the CLI application.
	// It serves as the entry point for the application and can have subcommands.
	var rootCmd = &cobra.Command{
		Use:   "taskctl",
		Short: "Task Scheduler CLI Client",
		Long:  "CLI Client for managing tasks in the Task Scheduler application.",
	}

	//Subcommands are added to the root command.
	rootCmd.AddCommand(commands.NewCreateCommand())
	rootCmd.AddCommand(commands.NewGetCommand())
	rootCmd.AddCommand(commands.NewListCommand())
	rootCmd.AddCommand(commands.NewUpdateCommand())
	//rootCmd.AddCommand(commands.NewDeleteCommand())
	//rootCmd.AddCommand(commands.NewHealthCheckCommand())
	//rootCmd.AddCommand(commands.NewScheduleCommand())

	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
