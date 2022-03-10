package cmd

import (
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <action|device|rule|state> -f FILENAME",
	Short: "Create a resource from file",
	Long:  `Create a resource from (JSON and YAML formats are accepted)`,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
