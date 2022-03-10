package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var actionsList []string
var assignWorkCmd = &cobra.Command{
	Use:   "work",
	Short: "Assign new work to device",
	Long: `Assign new work to device.

Assign a new work to device manually, by setting either a rule (takes precedence) or a list of actions

Examples:
  # Assign rule to device
  vaxctl assign work -d DEVICE_UID -r RULE_NAME
		
  # Assign actions to device
  vaxctl assign work -d DEVICE_UID -a "Press F1, Press F2"
	
  # Assign work to device from file
  vaxctl assign work -f work_assignment.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if filename == "" {
			if deviceUid == "" {
				fmt.Println("device UID must be set if no filename is given")
				cmd.Usage()
				os.Exit(2)
			}
			if name == "" && len(actionsList) == 0 {
				fmt.Println("either a rule or a list of actions must be set if no filename is given")
				cmd.Usage()
				os.Exit(2)
			}
		}
		err := model.AssignWork(deviceUid, name, actionsList, filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	assignCmd.AddCommand(assignWorkCmd)
	assignWorkCmd.Flags().StringVarP(&deviceUid, "device", "d", "", "uid of device")
	assignWorkCmd.RegisterFlagCompletionFunc("device", model.GetDeviceNamesForCompletion)
	assignWorkCmd.Flags().StringVarP(&name, "rule", "r", "", "name of rule to run")
	assignWorkCmd.RegisterFlagCompletionFunc("rule", model.GetRuleNamesForCompletion)
	assignWorkCmd.Flags().StringSliceVarP(&actionsList, "actions", "a", []string{}, "comma separated list of actions")
	assignWorkCmd.RegisterFlagCompletionFunc("actions", model.GetActionNamesForCompletion)
	assignWorkCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create the resource")
}
