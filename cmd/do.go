package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Complete the task",
	Run: func(cmd *cobra.Command, args []string) {
		nb, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal(err)
		}
		CompleteTask(nb)
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
