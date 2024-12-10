package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"task/db"
)

var completedCmd = &cobra.Command{
	Use:     "completed",
	Aliases: []string{"done"},
	Short:   "List completed tasks",
	Long:    "List all completed tasks during the day",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_, tasks, err := db.Completed()
		if err != nil {
			fmt.Println("Error listing tasks:", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found.")
			return
		}

		fmt.Println("Tasks:")
		for i, task := range tasks {
			fmt.Printf("%d. %s\n", i+1, task) // Display tasks with dynamic numbering
		}
	},
}

func init() {
	rootCmd.AddCommand(completedCmd)
}
