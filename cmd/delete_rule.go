package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var deleteRuleCmd = &cobra.Command{
	Use:   "rule",
	Short: "Delete rule by name or from file",
	Long: `Delete rule by name or from file.

JSON and YAML formats are accepted.
  
Examples:
  # delete rule from yaml
  vaxctl delete rule -f rule.yaml
    
  # create rule by name
  vaxctl delete rule -n RULE_NAME`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if (filename == "" && name == "") || (filename != "" && name != "") {
			fmt.Println("You must set either '-f' or '-n'")
			cmd.Usage()
			os.Exit(2)
		}
		err := model.DeleteResource("rule", filename, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteRuleCmd)

	deleteRuleCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to delete the resource")
	deleteRuleCmd.Flags().StringVarP(&name, "name", "n", "", "name of the resource to delete")
	deleteRuleCmd.RegisterFlagCompletionFunc("name", model.GetRuleNamesForCompletion)
}
