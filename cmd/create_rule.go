package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"
	"vaxctl/tui"

	"github.com/spf13/cobra"
)

var createRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Create rule from file",
	Long: `Create a new rule from file.

JSON and YAML formats are accepted.
  
Examples:
  # create rule from yaml
  vaxctl create rule -f rule.yaml

  # create rule from state in interactive mode
  vaxctl create rule -i STATE_ID`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if name != "" {
			err = tui.CreateRule(name)
		} else if filename != "" {
			err = model.CreateResource("rule", filename)
		} else {
			fmt.Println("You must set either a filename OR a state ID")
			cmd.Help()
			os.Exit(2)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	createCmd.AddCommand(createRuleCmd)
	createRuleCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create the resource")
	createRuleCmd.Flags().StringVarP(&name, "id", "i", "", "ID of the state (for interactive mode)")
	createRuleCmd.RegisterFlagCompletionFunc("id", model.GetStateIdsForCompletion)
}
