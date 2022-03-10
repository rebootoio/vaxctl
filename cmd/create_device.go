package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"
	"vaxctl/tui"

	"github.com/spf13/cobra"
)

var createDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Create device from file",
	Long: `Create a new device from file.

JSON and YAML formats are accepted.
  
Examples:
  # create device from json
  vaxctl create device -f device.json
    
  # create device from yaml
  vaxctl create device -f device.yaml

  # create device in interactive mode
  vaxctl create device -i`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if interactive {
			err = tui.CreateDevice()
		} else if filename != "" {
			err = model.CreateResource("device", filename)
		} else {
			fmt.Println("You must set either a filename OR enable interactive mode")
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
	createCmd.AddCommand(createDeviceCmd)
	createDeviceCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create the resource")
	createDeviceCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "open interactive mode")
}
