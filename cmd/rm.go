package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove a task",
	Run: func(cmd *cobra.Command, args []string) {
		nb, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		RemoveTask(nb)
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
