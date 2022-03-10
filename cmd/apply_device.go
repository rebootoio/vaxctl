package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var applyDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Create/Update device from file",
	Long: `Create/Update a new device from file.

JSON and YAML formats are accepted.
  
Examples:
  # apply device from json
  vaxctl apply device -f device.json
    
  # apply device from yaml
  vaxctl apply device -f device.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.ApplyResource("device", filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	applyCmd.AddCommand(applyDeviceCmd)
	applyDeviceCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create/update the resource")
	applyDeviceCmd.MarkFlagRequired("filename")
}
