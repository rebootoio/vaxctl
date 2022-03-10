package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var applyCredCmd = &cobra.Command{
	Use:   "cred",
	Short: "Create/Update cred from file",
	Long: `Create/Update a new cred from file.

JSON and YAML formats are accepted.
  
Examples:
  # apply cred from json
  vaxctl apply cred -f cred.json
    
  # apply cred from yaml
  vaxctl apply cred -f cred.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.ApplyResource("creds", filename)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	applyCmd.AddCommand(applyCredCmd)
	applyCredCmd.Flags().StringVarP(&filename, "filename", "f", "", "filename to use to create/update the resource")
	applyCredCmd.MarkFlagRequired("filename")
}
