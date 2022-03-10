package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"
	"vaxctl/tui"

	"github.com/spf13/cobra"
)

var createCredCmd = &cobra.Command{
	Use:   "cred",
	Short: "Create cred from file",
	Long: `Create a new cred from file.

JSON and YAML formats are accepted.
  
Examples:
  # create cred from json
  vaxctl create cred -f cred.json
    
  # create cred from yaml
  vaxctl create cred -f cred.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if interactive {
			err = tui.CreateCred()
		} else if filename != "" {
			err = model.CreateResource("creds", filename)
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
	createCmd.AddCommand(createCredCmd)
	createCredCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create the resource")
	createCredCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "open interactive mode")
}
