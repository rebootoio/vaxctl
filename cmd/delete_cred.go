package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var deleteCredCmd = &cobra.Command{
	Use:   "cred",
	Short: "Delete cred by id or from file",
	Long: `Delete cred by id or from file.

JSON and YAML formats are accepted.
  
Examples:
  # delete cred from yaml
  vaxctl delete cred -f cred.yaml
    
  # delete cred by name
  vaxctl delete cred -n NAME`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if (filename == "" && name == "") || (filename != "" && name != "") {
			fmt.Println("You must set either '-f' or '-n'")
			cmd.Usage()
			os.Exit(2)
		}
		err := model.DeleteResource("creds", filename, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteCredCmd)
	deleteCredCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to delete the resource")
	deleteCredCmd.Flags().StringVarP(&name, "name", "n", "", "name of the resource to delete")
	deleteCredCmd.RegisterFlagCompletionFunc("name", model.GetCredNamesForCompletion)
}
