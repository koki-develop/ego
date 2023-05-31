package cmd

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"
)

type egoOptions struct {
	NoNewline      bool
	EnableEscapes  bool
	DisableEscapes bool

	DisableStyle bool

	Foreground string
	Background string

	Bold          bool
	Faint         bool
	Italic        bool
	Underline     bool
	Blink         bool
	RapidBlink    bool
	Strikethrough bool

	Separator string

	Timestamp       bool
	TimestampFormat string
}

func ego(w io.Writer, args []string, options egoOptions) error {
	b := new(bytes.Buffer)

	if options.Timestamp {
		if _, err := b.WriteString(time.Now().Format(options.TimestampFormat) + " "); err != nil {
			return err
		}
	}

	if !options.DisableStyle {
		if err := style(b, options); err != nil {
			return err
		}
	}

	for i, arg := range args {
		if !options.DisableEscapes && options.EnableEscapes {
			if err := interpretEscapes(b, arg); err != nil {
				return err
			}
		} else {
			if _, err := b.WriteString(arg); err != nil {
				return err
			}
		}

		if i+1 != len(args) {
			if _, err := b.WriteString(options.Separator); err != nil {
				return err
			}
		}
	}

	if !options.DisableStyle {
		if err := resetStyle(b); err != nil {
			return err
		}
	}

	if !options.NoNewline {
		if _, err := b.WriteRune('\n'); err != nil {
			return err
		}
	}

	if _, err := io.Copy(w, b); err != nil {
		return err
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
	ansiFg = map[string]string{
		"black":     "30",
		"red":       "31",
		"green":     "32",
		"yellow":    "33",
		"blue":      "34",
		"magenta":   "35",
		"cyan":      "36",
		"white":     "37",
		"hiblack":   "90",
		"hired":     "91",
		"higreen":   "92",
		"hiyellow":  "93",
		"hiblue":    "94",
		"himagenta": "95",
		"hicyan":    "96",
		"hiwhite":   "97",
	}

	ansiBg = map[string]string{
		"black":     "40",
		"red":       "41",
		"green":     "42",
		"yellow":    "43",
		"blue":      "44",
		"magenta":   "45",
		"cyan":      "46",
		"white":     "47",
		"hiblack":   "100",
		"hired":     "101",
		"higreen":   "102",
		"hiyellow":  "103",
		"hiblue":    "104",
		"himagenta": "105",
		"hicyan":    "106",
		"hiwhite":   "107",
	}
)

func style(w io.Writer, options egoOptions) error {
	if options.Foreground != "" {
		ansi, ok := ansiFg[options.Foreground]
		if !ok {
			return fmt.Errorf("unsupported foreground color: %s", options.Foreground)
		}
		if err := writeAnsi(w, ansi); err != nil {
			return err
		}
	}

	if options.Background != "" {
		ansi, ok := ansiBg[options.Background]
		if !ok {
			return fmt.Errorf("unsupported background color: %s", options.Background)
		}
		if err := writeAnsi(w, ansi); err != nil {
			return err
		}
	}

	if options.Bold {
		if err := writeAnsi(w, "1"); err != nil {
			return err
		}
	}

	if options.Faint {
		if err := writeAnsi(w, "2"); err != nil {
			return err
		}
	}

	if options.Italic {
		if err := writeAnsi(w, "3"); err != nil {
			return err
		}
	}

	if options.Underline {
		if err := writeAnsi(w, "4"); err != nil {
			return err
		}
	}

	if options.Blink {
		if err := writeAnsi(w, "5"); err != nil {
			return err
		}
	}

	if options.RapidBlink {
		if err := writeAnsi(w, "6"); err != nil {
			return err
		}
	}

	if options.Strikethrough {
		if err := writeAnsi(w, "9"); err != nil {
			return err
		}
	}

	return nil
}

func resetStyle(w io.Writer) error {
	if err := writeAnsi(w, "0"); err != nil {
		return err
	}

	return nil
}

func writeAnsi(w io.Writer, ansi string) error {
	if _, err := w.Write([]byte{'\x1b', '['}); err != nil {
		return err
	}
	if _, err := w.Write([]byte(ansi)); err != nil {
		return err
	}
	if _, err := w.Write([]byte{'m'}); err != nil {
		return err
	}

	return nil
}
