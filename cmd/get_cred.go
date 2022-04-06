package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var getCredCmd = &cobra.Command{
	Use:   "cred",
	Short: "Get one or many credentials",
	Long: `Get credentials details.

Prints a table of the most important information about the credentials

Examples:
  # List all credentials
  vaxctl get cred
		
  # Get cred by name as yaml
  vaxctl get cred -n NAME -o yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.PrintCreds(name, output)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getCredCmd)
	getCredCmd.Flags().StringVarP(&name, "name", "n", "", "name of resource (if not set all are returned)")
	getCredCmd.RegisterFlagCompletionFunc("name", model.GetCredNamesForCompletion)
	getCredCmd.Flags().StringVarP(&output, "output", "o", "", "output format (default is table). One of: json|yaml")
}
