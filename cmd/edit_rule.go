package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"
	"vaxctl/tui"

	"github.com/spf13/cobra"
)

var editRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "edit rule from server",
	Long: `edit a rule in interactive mode.

Can be exported to JSON and YAML formats.
  
Examples:
  # edit rule from server
  vaxctl edit rule -n RULE_NAME`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := tui.EditRule(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	editCmd.AddCommand(editRuleCmd)
	editRuleCmd.Flags().StringVarP(&name, "name", "n", "", "name of the rule to edit")
	editRuleCmd.MarkFlagRequired("name")
	editRuleCmd.RegisterFlagCompletionFunc("name", model.GetRuleNamesForCompletion)
}
