package cmd

import (
	"github.com/spf13/cobra"
)

var mandatoryFlag, commentsFlag bool

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new resource template",
	Long:  `Generate a local file with the mandatory/optional fields for a resource`,
	Args:  cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "output file (if not set with print to stdout)")
	generateCmd.PersistentFlags().BoolVarP(&mandatoryFlag, "mandatory", "m", false, "print only mandatory fields")
	generateCmd.PersistentFlags().BoolVarP(&commentsFlag, "comments", "c", false, "print with comments")
}
