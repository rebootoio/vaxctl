package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var getRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Get one or many rules",
	Long: `Get rule details.

Prints a table of the most important information about the rules

Examples:
  # List all rules
  vaxctl get rule
		
  # Get rule by name as yaml
  vaxctl get rule -n RULE_NAME -o yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.PrintRules(name, verbose, output)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getRuleCmd)
	getRuleCmd.Flags().StringVarP(&name, "name", "n", "", "name of resource (if not set all are returned)")
	getRuleCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "include screenshot & ocr_text values in yaml/json")
	getRuleCmd.RegisterFlagCompletionFunc("name", model.GetRuleNamesForCompletion)
	getRuleCmd.Flags().StringVarP(&output, "output", "o", "", "output format (default is table). One of: json|yaml")
}
