package cmd

import (
	"fmt"
	"os"
	"vaxctl/model"

	"github.com/spf13/cobra"
)

var generateCredCmd = &cobra.Command{
	Use:   "cred",
	Short: "Generate a cred template",
	Long: `Get a new credential.

Examples:
  # Generate a cred template and print to screen
  vaxctl generate cred
		
  # Get cred template in a file
  vaxctl generate cred -f new_cred.yaml`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := model.GenerateCred(filename, mandatoryFlag, commentsFlag)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	generateCmd.AddCommand(generateCredCmd)
}
