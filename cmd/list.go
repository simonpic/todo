package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Todo tasks",
	Run: func(cmd *cobra.Command, args []string) {
		for i, task := range ReadTasks() {
			var status string
			if task.Completed {
				status = "X"
			} else {
				status = " "
			}
			fmt.Printf("%s %d - %s\n", status, i+1, task.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
