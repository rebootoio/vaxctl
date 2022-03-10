package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var applyActionCmd = &cobra.Command{
	Use:   "action",
	Short: "Create/Update action from file",
	Long: `Create/Update a new action from file.

JSON and YAML formats are accepted.
  
Examples:
  # apply action from json
  vaxctl apply action -f action.json
    
  # apply action from yaml
  vaxctl apply action -f action.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.ApplyResource("action", filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	applyCmd.AddCommand(applyActionCmd)
	applyActionCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create/update the resource")
	applyActionCmd.MarkFlagRequired("filename")
}
