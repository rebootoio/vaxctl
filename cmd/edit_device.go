package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"
	"vaxctl/tui"

	"github.com/spf13/cobra"
)

var editDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "edit device from server",
	Long: `edit a device in interactive mode.

Can be exported to JSON and YAML formats.
  
Examples:
  # edit device from server
  vaxctl edit device -n DEVICE_UID`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := model.GetDevices(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = tui.EditDevice(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	editCmd.AddCommand(editDeviceCmd)
	editDeviceCmd.Flags().StringVarP(&name, "name", "n", "", "UID of the device to edit")
	editDeviceCmd.MarkFlagRequired("name")
	editDeviceCmd.RegisterFlagCompletionFunc("name", model.GetDeviceNamesForCompletion)
}
