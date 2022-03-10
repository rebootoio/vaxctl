package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"
	"vaxctl/tui"

	"github.com/spf13/cobra"
)

var editActionCmd = &cobra.Command{
	Use:   "action",
	Short: "edit action from server",
	Long: `edit an action in interactive mode.

Can be exported to JSON and YAML formats.
  
Examples:
  # edit action from server
  vaxctl edit action -n ACTION_NAME`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := tui.EditAction(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	editCmd.AddCommand(editActionCmd)
	editActionCmd.Flags().StringVarP(&name, "name", "n", "", "name of the action to edit")
	editActionCmd.MarkFlagRequired("name")
	editActionCmd.RegisterFlagCompletionFunc("name", model.GetActionNamesForCompletion)
}
