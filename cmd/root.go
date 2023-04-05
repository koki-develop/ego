package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var options egoOptions

var rootCmd = &cobra.Command{
	Use: "ego",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := ego(os.Stdout, args, options); err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&(options.NoNewline), "no-newline", "n", false, "do not print the trailing newline character")
}
