package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var deleteActionCmd = &cobra.Command{
	Use:   "action",
	Short: "Delete action by name or from file",
	Long: `Delete action by name or from file.

JSON and YAML formats are accepted.
  
Examples:
  # delete action from yaml
  vaxctl delete action -f action.yaml
    
  # create action by name
  vaxctl delete action -n ACTION_NAME`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if (filename == "" && name == "") || (filename != "" && name != "") {
			fmt.Println("You must set either '-f' or '-n'")
			cmd.Usage()
			os.Exit(2)
		}
		err := model.DeleteResource("action", filename, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteActionCmd)

	deleteActionCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to delete the resource")
	deleteActionCmd.Flags().StringVarP(&name, "name", "n", "", "name of the resource to delete")
	deleteActionCmd.RegisterFlagCompletionFunc("name", model.GetActionNamesForCompletion)
}
