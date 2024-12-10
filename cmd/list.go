package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"task/db"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l"},
	Short:   "List tasks",
	Long:    "List all active tasks on your list",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_, tasks, err := db.List(db.TaskBucket)
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
	rootCmd.AddCommand(listCmd)
}
