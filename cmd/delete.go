package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <action|device|rule|state> [flags]",
	Short: "Delete a resource",
	Long: `Delete resources by filenames or names (JSON and YAML formats are accepted).

Only one type of the arguments may be specified: filenames or names.`,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
