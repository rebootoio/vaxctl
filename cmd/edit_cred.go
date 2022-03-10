package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"
	"vaxctl/tui"

	"github.com/spf13/cobra"
)

var editCredCmd = &cobra.Command{
	Use:   "cred",
	Short: "edit cred from server",
	Long: `edit an cred in interactive mode.

Can be exported to JSON and YAML formats.
  
Examples:
  # edit cred from server
  vaxctl edit cred -n CRED_NAME`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := tui.EditCred(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	editCmd.AddCommand(editCredCmd)
	editCredCmd.Flags().StringVarP(&name, "name", "n", "", "name of the credential to edit")
	editCredCmd.MarkFlagRequired("name")
	editCredCmd.RegisterFlagCompletionFunc("name", model.GetCredNamesForCompletion)
}
