package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var getStateCmd = &cobra.Command{
	Use:   "state",
	Short: "Get one or many states",
	Long: `Get state details.

Prints a table of the most important information about the states.

Examples:
  # List all states
  vaxctl get state
		
  # Get state by id as yaml
  vaxctl get state -i STATE_ID -o yaml
	
  # Get open states
  vaxctl get state -t open
	
  # Get resolved states as yaml
  vaxctl get state -t resolved -o yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.PrintStates(name, filename, deviceUid, regex, verbose, output)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	getCmd.AddCommand(getStateCmd)
	getStateCmd.Flags().StringVarP(&name, "id", "i", "", "name of resource (if not set all are returned)")
	getStateCmd.RegisterFlagCompletionFunc("id", model.GetStateIdsForCompletion)
	getStateCmd.Flags().StringVarP(&filename, "type", "t", "", "type of states (allowed values are: open, unknown, resolved. If not set all are returned)")
	getStateCmd.RegisterFlagCompletionFunc("type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"open", "unknown", "resolved"}, cobra.ShellCompDirectiveNoFileComp
	})
	getStateCmd.Flags().StringVarP(&deviceUid, "device", "d", "", "get states for specific device (if not set all are returned)")
	getStateCmd.RegisterFlagCompletionFunc("device", model.GetDeviceNamesForCompletion)
	getStateCmd.Flags().StringVarP(&regex, "regex", "r", "", "get states by matching regex (if not set all are returned)")
	getStateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show full OCR text")
}
