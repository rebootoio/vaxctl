package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var generateStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Generate a state template",
	Long: `Get a new state.

Examples:
  # Generate a state template and print to screen
  vaxctl generate state
		
  # Get state template in a file
  vaxctl generate state -f new_state.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.GenerateState(filename, mandatoryFlag, commentsFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	generateCmd.AddCommand(generateStateCmd)
}
