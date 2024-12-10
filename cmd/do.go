package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"task/db"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Complete a task",
	Long:  "Mark a task as completed, remove from the list and save it in the list of completed tasks",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		numTask, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: Invalid task ID.")
			return
		}
		err = db.Do(numTask)
		if err != nil {
			fmt.Println("Internal server error")
			return
		} else {
			fmt.Println("Task marked as completed! ")
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
