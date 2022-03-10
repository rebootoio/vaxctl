package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var applyStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Create/Update state from file",
	Long: `Create/Update a new state from file.

JSON and YAML formats are accepted.
  
Examples:
  # apply state from json
  vaxctl apply state -f state.json
    
  # apply state from yaml
  vaxctl apply state -f state.yaml`,
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
	applyCmd.AddCommand(applyStateCmd)
	applyStateCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create/update the resource")
	applyStateCmd.MarkFlagRequired("filename")
}
