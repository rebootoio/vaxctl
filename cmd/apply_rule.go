package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var applyRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Create/Update rule from file",
	Long: `Create/Update a new rule from file.

JSON and YAML formats are accepted.
  
Examples:
  # apply rule from json
  vaxctl apply rule -f rule.json
    
  # apply rule from yaml
  vaxctl apply rule -f rule.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.ApplyResource("rule", filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	applyCmd.AddCommand(applyRuleCmd)
	applyRuleCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create/update the resource")
	applyRuleCmd.MarkFlagRequired("filename")
}
