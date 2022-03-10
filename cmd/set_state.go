package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var setStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Set state as resolved",
	Long: `Set state status to be resolved.

Manually set a device state to resolved

Examples:
  # Set state as resolved for device 
  vaxctl set state -d DEVICE_UID`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.SetStateAsResolved(name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	setCmd.AddCommand(setStateCmd)
	setStateCmd.Flags().StringVarP(&name, "device", "d", "", "uid of device (if not set all are returned)")
	setStateCmd.RegisterFlagCompletionFunc("device", model.GetDeviceNamesForCompletion)
	setStateCmd.MarkFlagRequired("device")
}
