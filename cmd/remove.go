package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
	"task/db"
)

var shouldRemoveAll bool
var removeCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete a task",
	Long:  "Delete a task from your list",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error: Invalid task ID.")
			return
		}
		ids, tasks, err := db.List(db.TaskBucket)
		if err != nil {
			fmt.Println("Error fetching tasks:", err)
			return
		}
		if id > len(tasks) {
			fmt.Println("Error: Task number out of range.")
			return
		}
		dbid := ids[id-1]
		if err := db.Remove(dbid, db.TaskBucket, shouldRemoveAll); err != nil {
			fmt.Println("Error deleting task:", err)
		} else {
			fmt.Println("Task deleted successfully! ")
		}
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&shouldRemoveAll, "all", "c", false, "Remove all completed tasks")
	rootCmd.AddCommand(removeCmd)
}
