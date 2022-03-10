package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var deleteDeviceCmd = &cobra.Command{
	Use:   "device",
	Short: "Delete device by name or from file",
	Long: `Delete device by name or from file.

JSON and YAML formats are accepted.
  
Examples:
  # delete device from yaml
  vaxctl delete device -f device.yaml
    
  # create device by name
  vaxctl delete device -n DEVICE_NAME`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if (filename == "" && name == "") || (filename != "" && name != "") {
			fmt.Println("You must set either '-f' or '-n'")
			cmd.Usage()
			os.Exit(2)
		}
		err := model.DeleteResource("device", filename, name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	deleteCmd.AddCommand(deleteDeviceCmd)
	deleteDeviceCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to delete the resource")
	deleteDeviceCmd.Flags().StringVarP(&name, "uid", "u", "", "uid of the resource to delete")
	deleteDeviceCmd.RegisterFlagCompletionFunc("uid", model.GetDeviceNamesForCompletion)
}
