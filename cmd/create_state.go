package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var createStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Create state from file",
	Long: `Create a new state from file.

JSON and YAML formats are accepted.
  
Examples:
  # create state from json
  vaxctl create state -f state.json
    
  # create state from yaml
  vaxctl create state -f state.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.CreateOrUpdateState(filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	createCmd.AddCommand(createStateCmd)
	createStateCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create the resource")
	createStateCmd.MarkFlagRequired("filename")
}
