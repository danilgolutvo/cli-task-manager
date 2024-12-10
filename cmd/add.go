package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
	"task/db"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"addition"},
	Short:   "Add a task",
	Long:    "Save your tasks",
	Args:    cobra.RangeArgs(1, 100),
	RunE: func(cmd *cobra.Command, args []string) error {
		task := strings.Join(args, " ")
		err := db.Add(task, db.TaskBucket)
		if err != nil {
			return err
		}
		fmt.Printf("Task %s successfully added to your list\n", task)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
