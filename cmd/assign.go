package cmd

import (
	"github.com/spf13/cobra"
)

var assignCmd = &cobra.Command{
	Use:   "assign <resource> [FLAGS]",
	Short: "Assign work",
	Long:  `Assign work to resource`,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(assignCmd)
}
