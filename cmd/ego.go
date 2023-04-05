package cmd

import (
	"io"
)

type egoOptions struct {
	NoNewline bool
}

func ego(w io.Writer, args []string, options egoOptions) error {
	for i, arg := range args {
		if _, err := w.Write([]byte(arg)); err != nil {
			return err
		}

		if i+1 != len(args) {
			if _, err := w.Write([]byte(" ")); err != nil {
				return err
			}
		}
	}
	if !options.NoNewline {
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}
