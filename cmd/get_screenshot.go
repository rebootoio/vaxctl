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

Saves the screenshot of a state - either by ID or device UID (will bring the latest for the device)

Examples:
  # Get screenshot for state ID
  vaxctl get screenshot -i STATE_ID -f screenshot.png
		
  # Get screenshot for latest state of device
  vaxctl get screenshot -d DEVICE_UID -f screenshot.png`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if name == "" && deviceUid == "" {
			fmt.Println("You must set either a device UID OR a state ID")
			cmd.Usage()
			os.Exit(2)
		}
		err := model.GetScreenshot(name, deviceUid, filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getScreenshotCmd)
	getScreenshotCmd.Flags().StringVarP(&name, "id", "i", "", "id of the state")
	getScreenshotCmd.RegisterFlagCompletionFunc("id", model.GetStateIdsForCompletion)
	getScreenshotCmd.Flags().StringVarP(&filename, "filename", "f", "", "output file for screenshot")
	getScreenshotCmd.Flags().StringVarP(&deviceUid, "device", "d", "", "get states for specific device (if not set all are returned)")
	getScreenshotCmd.RegisterFlagCompletionFunc("device", model.GetDeviceNamesForCompletion)
	getScreenshotCmd.MarkFlagRequired("filename")
}
