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

	rootCmd.Flags().StringVar(&(options.Foreground), "foreground", "", "foreground color")
	rootCmd.Flags().StringVar(&(options.Background), "background", "", "background color")
	rootCmd.Flags().BoolVar(&(options.Bold), "bold", false, "bold strings")
	rootCmd.Flags().BoolVar(&(options.Faint), "faint", false, "faint strings")
	rootCmd.Flags().BoolVar(&(options.Italic), "italic", false, "italicize strings")
	rootCmd.Flags().BoolVar(&(options.Underline), "underline", false, "underline strings")
	rootCmd.Flags().BoolVar(&(options.Blink), "blink", false, "blink strings")
	rootCmd.Flags().BoolVar(&(options.RapidBlink), "rapid-blink", false, "rapid blink strings")
	rootCmd.Flags().BoolVar(&(options.Strikethrough), "strikethrough", false, "strikethrough strings")
}
