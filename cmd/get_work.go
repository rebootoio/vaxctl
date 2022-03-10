package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var showDetails bool
var latest bool

var getWorkCmd = &cobra.Command{
	Use:   "work",
	Short: "Get one or many work",
	Long: `Get work details.

Prints a table of the most important information about the works

Examples:
  # List all works
  vaxctl get work
		
  # Get latest work with details by device
  vaxctl get work -d DEVICE_UID -v -l`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.GetWorks(filename, name, showDetails, latest, output)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getWorkCmd)
	getWorkCmd.Flags().StringVarP(&name, "device", "d", "", "uid of device (if not set all are returned)")
	getWorkCmd.RegisterFlagCompletionFunc("device", model.GetDeviceNamesForCompletion)
	getWorkCmd.Flags().BoolVarP(&showDetails, "verbose", "v", false, "show verbose view of running/completed work")
	getWorkCmd.Flags().BoolVarP(&latest, "latest", "l", false, "show only latest work")
	getWorkCmd.Flags().StringVarP(&filename, "id", "i", "", "id of resource (if not set all are returned)")
	getWorkCmd.RegisterFlagCompletionFunc("id", model.GetWorkIdsForCompletion)
}
