package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var getActionCmd = &cobra.Command{
	Use:   "action",
	Short: "Get one or many actions",
	Long: `Get action details.

Prints a table of the most important information about the actions

Examples:
  # List all actions
  vaxctl get action
		
  # Get action by name as yaml
  vaxctl get action -n ACTION_NAME -o yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.PrintActions(name, output)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getActionCmd)
	getActionCmd.Flags().StringVarP(&name, "name", "n", "", "name of resource (if not set all are returned)")
	getActionCmd.RegisterFlagCompletionFunc("name", model.GetActionNamesForCompletion)
	getActionCmd.Flags().StringVarP(&output, "output", "o", "", "output format (default is table). One of: json|yaml")
}
