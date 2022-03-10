package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var getDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Get one or many devices",
	Long: `Get device details.

Prints a table of the most important information about the devices

Examples:
  # List all devices
  vaxctl get device
		
  # Get device by uid as yaml
  vaxctl get device -n UID -o yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.PrintDevices(name, output)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getDeviceCmd)
	getDeviceCmd.Flags().StringVarP(&name, "uid", "n", "", "uid of resource (if not set all are returned)")
	getDeviceCmd.RegisterFlagCompletionFunc("uid", model.GetDeviceNamesForCompletion)
}
