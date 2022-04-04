package cmd

import (
	"fmt"
	"os"
	"strconv"
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
			_, err = model.GetStates(name, "", "", "")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = tui.CreateRule(name)
		} else if output != "" {
			states, err := model.GetStates("", "", output, "")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			numberOfStates := len(states)
			if numberOfStates == 0 {
				fmt.Printf("Error: State for UID '%s' was not found\n", output)
				os.Exit(1)
			} else {
				lastState := states[numberOfStates-1]
				if lastState.Resolved {
					fmt.Printf("Error: Open State for UID '%s' was not found, last resolved state ID is '%d'\n", output, lastState.StateId)
					os.Exit(1)
				} else {
					err = tui.CreateRule(strconv.Itoa(lastState.StateId))
				}
			}
		} else if filename != "" {
			err = model.CreateResource("rule", filename)
		} else {
			fmt.Println("You must set either a filename OR a state ID OR a device UID")
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
	createRuleCmd.Flags().StringVarP(&output, "device", "d", "", "device UID of the open state (for interactive mode)")
	createRuleCmd.RegisterFlagCompletionFunc("id", model.GetStateIdsForCompletion)
	createRuleCmd.RegisterFlagCompletionFunc("device", model.GetDeviceNamesForCompletion)
}
