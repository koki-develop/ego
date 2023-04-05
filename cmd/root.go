package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var options egoOptions

var rootCmd = &cobra.Command{
	Use:          "ego [flags] [strings]",
	Short:        "echo alternative written in Go",
	Long:         "echo alternative written in Go.",
	SilenceUsage: true,
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
	rootCmd.Flags().SortFlags = false

	rootCmd.Flags().BoolVarP(&(options.NoNewline), "no-newline", "n", false, "do not print the trailing newline character")
	rootCmd.Flags().BoolVarP(&(options.EnableEscapes), "enable-escapes", "e", true, "enable interpretation of backslash escapes")

	rootCmd.Flags().StringVar(&(options.Foreground), "fg", "", "foreground color")
	rootCmd.Flags().StringVar(&(options.Background), "bg", "", "background color")
}
