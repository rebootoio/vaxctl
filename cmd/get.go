package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display one or many resources",
	Long:  `Prints a table of the most important information about the specified resources`,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(getCmd)
}
