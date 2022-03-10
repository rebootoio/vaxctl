package cmd

import (
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply <action|device|rule|state> -f FILENAME",
	Short: "Create/Update a resource from file",
	Long:  `Create/Update a resource from a file (JSON and YAML formats are accepted)`,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
