package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"task/db"
)

const DBPath = "tasks.db"

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a cli tool for performing basic managing operations with tasks",
	Long:  "task is a cli tool for performing basic managing operations with tasks like adding, listing, removing, marking.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}

func init() {
	db.Init(DBPath)
}
