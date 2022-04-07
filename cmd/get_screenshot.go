package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var getScreenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Get screenshot file from state",
	Long: `Get screenshot of a specific state.

Get the screenshot of:
* state - either by ID or device UID (will bring the latest for the device)
* rule - by rule name

Examples:
  # Get screenshot for state ID
  vaxctl get screenshot -i STATE_ID -f screenshot.png
		
  # Get screenshot for latest state of device
  vaxctl get screenshot -d DEVICE_UID -f screenshot.png

  # Get screenshot for rule
  vaxctl get screenshot -r RULE_NAME -f screenshot.png`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" && deviceUid == "" && output == "" {
			fmt.Println("You must set either a device UID OR a state ID OR a rule name")
			cmd.Usage()
			os.Exit(2)
		}
		err := model.GetScreenshot(name, deviceUid, output, filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getScreenshotCmd)
	getScreenshotCmd.Flags().StringVarP(&name, "id", "i", "", "id of the state")
	getScreenshotCmd.Flags().StringVarP(&deviceUid, "device", "d", "", "uid of device")
	getScreenshotCmd.Flags().StringVarP(&output, "rule", "r", "", "rule name")
	getScreenshotCmd.Flags().StringVarP(&filename, "filename", "f", "", "output file for screenshot (if non provided it will open it)")
	getScreenshotCmd.RegisterFlagCompletionFunc("id", model.GetStateIdsForCompletion)
	getScreenshotCmd.RegisterFlagCompletionFunc("device", model.GetDeviceNamesForCompletion)
	getScreenshotCmd.RegisterFlagCompletionFunc("rule", model.GetRuleNamesForCompletion)
}
