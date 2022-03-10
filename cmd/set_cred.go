package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var setCredCmd = &cobra.Command{
	Use:   "cred",
	Short: "Set credential as default",
	Long: `Set credential to be default.

Set a credential to be the default

Examples:
  # Set credential as default 
  vaxctl set cred -n CRED_NAME`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.SetCredsAsDefault(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	setCmd.AddCommand(setCredCmd)
	setCredCmd.Flags().StringVarP(&name, "name", "n", "", "name of the credentials")
	setCredCmd.RegisterFlagCompletionFunc("name", model.GetCredNamesForCompletion)
	setCredCmd.MarkFlagRequired("name")
}
