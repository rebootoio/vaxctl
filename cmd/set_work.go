package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var setWorkCmd = &cobra.Command{
	Use:   "work",
	Short: "Set work status",
	Long: `Set work status.

Manually complete a pending work and set its status

Examples:
  # Set work status to failed for device 
  vaxctl set work -d DEVICE_UID -s failure`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if filename != "success" && filename != "failure" {
			fmt.Println("Status is only allowed to be 'success' or 'failure")
			cmd.Usage()
			os.Exit(2)
		}
		err := model.SetWork(name, filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	setCmd.AddCommand(setWorkCmd)
	setWorkCmd.Flags().StringVarP(&name, "device", "d", "", "uid of device (if not set all are returned)")
	setWorkCmd.RegisterFlagCompletionFunc("device", model.GetDeviceNamesForCompletion)
	setWorkCmd.MarkFlagRequired("device")
	setWorkCmd.Flags().StringVarP(&filename, "status", "s", "", "status to set (allowed values are: success & failure)")
	setWorkCmd.MarkFlagRequired("status")
	setWorkCmd.RegisterFlagCompletionFunc("status", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"success", "failure"}, cobra.ShellCompDirectiveNoFileComp
	})

}
