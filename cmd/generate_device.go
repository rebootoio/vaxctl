package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var generateDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Generate a device template",
	Long: `Get a new device.

Examples:
  # Generate a device template and print to screen
  vaxctl generate device
		
  # Get device template in a file
  vaxctl generate device -f new_device.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.GenerateDevice(filename, mandatoryFlag, commentsFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	generateCmd.AddCommand(generateDeviceCmd)
}
