package cmd

import (
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <action|device|rule> -n NAME",
	Short: "Edit a resource",
	Long:  `Edit a resource (open an interactive shell to edit a resource from the server)`,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(editCmd)
}
