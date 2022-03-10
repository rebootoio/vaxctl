package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var generateRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Generate a rule template",
	Long: `Get a new rule.

Examples:
  # Generate a rule template and print to screen
  vaxctl generate rule
		
  # Get rule template in a file
  vaxctl generate rule -f new_rule.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.GenerateRule(filename, mandatoryFlag, commentsFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	generateCmd.AddCommand(generateRuleCmd)
}
