package cmd

import (
	"fmt"
	"io"
	"strconv"
)

type egoOptions struct {
	NoNewline     bool
	EnableEscapes bool

	Foreground string
	Background string
}

func ego(w io.Writer, args []string, options egoOptions) error {
	if err := style(w, options); err != nil {
		return err
	}

	for i, arg := range args {
		if options.EnableEscapes {
			if err := interpretEscapes(w, arg); err != nil {
				return err
			}
		} else {
			if _, err := w.Write([]byte(arg)); err != nil {
				return err
			}
		}

		if i+1 != len(args) {
			if _, err := w.Write([]byte{' '}); err != nil {
				return err
			}
		}
	}

	if err := resetStyle(w); err != nil {
		return err
	}

	if !options.NoNewline {
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	return nil
}

func interpretEscapes(w io.Writer, s string) error {
	escaped := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if escaped {
			var err error
			switch c {
			case 'a':
				_, err = w.Write([]byte{'\a'})
			case 'b':
				_, err = w.Write([]byte{'\b'})
			case 'e':
				_, err = w.Write([]byte{'\x1b'})
			case 'f':
				_, err = w.Write([]byte{'\f'})
			case 'n':
				_, err = w.Write([]byte{'\n'})
			case 't':
				_, err = w.Write([]byte{'\t'})
			case 'v':
				_, err = w.Write([]byte{'\v'})
			case 'r':
				_, err = w.Write([]byte{'\r'})
			case '\\':
				_, err = w.Write([]byte{'\\'})
			case '"':
				_, err = w.Write([]byte{'"'})
			case '\'':
				_, err = w.Write([]byte{'\''})
			case '0':
				_, err = w.Write([]byte{'\x00'})
			case 'x':
				if i+2 < len(s) {
					num, err := strconv.ParseInt(s[i+1:i+3], 16, 8)
					if err == nil {
						if _, err := w.Write([]byte{byte(num)}); err != nil {
							return err
						}
						i += 2
					} else {
						if _, err := w.Write([]byte{'\\', c}); err != nil {
							return err
						}
					}
				} else {
					_, err = w.Write([]byte{'\\', c})
				}
			case 'o':
				if i+3 < len(s) {
					num, err := strconv.ParseInt(s[i+1:i+4], 8, 8)
					if err == nil {
						if _, err := w.Write([]byte{byte(num)}); err != nil {
							return err
						}
						i += 3
					} else {
						if _, err := w.Write([]byte{'\\', c}); err != nil {
							return err
						}
					}
				} else {
					_, err = w.Write([]byte{'\\', c})
				}
			default:
				_, err = w.Write([]byte{'\\', c})
			}
			if err != nil {
				return err
			}
			escaped = false
		} else if c == '\\' {
			escaped = true
		} else {
			_, err := w.Write([]byte{c})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

var (
	ansiFg = map[string][]byte{
		"black":     []byte("[30m"),
		"red":       []byte("[31m"),
		"green":     []byte("[32m"),
		"yellow":    []byte("[33m"),
		"blue":      []byte("[34m"),
		"magenta":   []byte("[35m"),
		"cyan":      []byte("[36m"),
		"white":     []byte("[37m"),
		"hiblack":   []byte("[90m"),
		"hired":     []byte("[91m"),
		"higreen":   []byte("[92m"),
		"hiyellow":  []byte("[93m"),
		"hiblue":    []byte("[94m"),
		"himagenta": []byte("[95m"),
		"hicyan":    []byte("[96m"),
		"hiwhite":   []byte("[97m"),
	}

	ansiBg = map[string][]byte{
		"black":     []byte("[40m"),
		"red":       []byte("[41m"),
		"green":     []byte("[42m"),
		"yellow":    []byte("[43m"),
		"blue":      []byte("[44m"),
		"magenta":   []byte("[45m"),
		"cyan":      []byte("[46m"),
		"white":     []byte("[47m"),
		"hiblack":   []byte("[100m"),
		"hired":     []byte("[101m"),
		"higreen":   []byte("[102m"),
		"hiyellow":  []byte("[103m"),
		"hiblue":    []byte("[104m"),
		"himagenta": []byte("[105m"),
		"hicyan":    []byte("[106m"),
		"hiwhite":   []byte("[107m"),
	}
)

func style(w io.Writer, options egoOptions) error {
	if options.Foreground != "" {
		ansi, ok := ansiFg[options.Foreground]
		if !ok {
			return fmt.Errorf("unsupported foreground color: %s", options.Foreground)
		}
		if _, err := w.Write([]byte{'\x1b'}); err != nil {
			return err
		}
		if _, err := w.Write(ansi); err != nil {
			return err
		}
	}

	if options.Background != "" {
		ansi, ok := ansiBg[options.Background]
		if !ok {
			return fmt.Errorf("unsupported background color: %s", options.Background)
		}
		if _, err := w.Write([]byte{'\x1b'}); err != nil {
			return err
		}
		if _, err := w.Write(ansi); err != nil {
			return err
		}
	}

	return nil
}

func resetStyle(w io.Writer) error {
	if _, err := w.Write([]byte{'\x1b'}); err != nil {
		return err
	}
	if _, err := w.Write([]byte("[0m")); err != nil {
		return err
	}

	return nil
}
