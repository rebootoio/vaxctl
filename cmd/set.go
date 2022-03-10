package cmd

import (
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Change state of a resource",
	Long:  `Manually chage the current status of a resource`,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(setCmd)
}
