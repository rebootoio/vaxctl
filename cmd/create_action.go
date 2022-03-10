package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"
	"vaxctl/tui"

	"github.com/spf13/cobra"
)

var listTypes, verbose bool

var createActionCmd = &cobra.Command{
	Use:   "action",
	Short: "Create action from file",
	Long: `Create a new action from file.

JSON and YAML formats are accepted.
  
Examples:
  # create action from json
  vaxctl create action -f action.json
    
  # create action from yaml
  vaxctl create action -f action.yaml

  # create action in interactive mode
  vaxctl create action -i`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if listTypes {
			if verbose {
				err = model.ListDetailedActionTypes()
			} else {
				err = model.ListActionTypes()
			}
		} else {
			if interactive {
				err = tui.CreateAction()
			} else if filename != "" {
				err = model.CreateResource("action", filename)
			} else {
				fmt.Println("You must set either a filename OR enable interactive mode")
				cmd.Help()
				os.Exit(2)
			}
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	createCmd.AddCommand(createActionCmd)
	createActionCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create the resource")
	createActionCmd.Flags().BoolVarP(&listTypes, "list-types", "I", false, "list available action types")
	createActionCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show detailed data for each action type")
	createActionCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "open interactive mode")
}
