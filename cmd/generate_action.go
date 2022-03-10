package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var generateActionCmd = &cobra.Command{
	Use:   "action",
	Short: "Generate a action template",
	Long: `Get a new action.

Examples:
  # Generate a action template and print to screen
  vaxctl generate action
		
  # Get action template in a file
  vaxctl generate action -f new_action.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.GenerateAction(filename, mandatoryFlag, commentsFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	generateCmd.AddCommand(generateActionCmd)
}
