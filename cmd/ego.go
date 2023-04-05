package cmd

import (
	"io"
	"strconv"
	"strings"
)

type egoOptions struct {
	NoNewline     bool
	EnableEscapes bool
}

func ego(w io.Writer, args []string, options egoOptions) error {
	for i, arg := range args {
		if options.EnableEscapes {
			arg = interpretEscapes(arg)
		}

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

func interpretEscapes(s string) string {
	var b strings.Builder

	escaped := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if escaped {
			switch c {
			case 'a':
				b.WriteRune('\a')
			case 'b':
				b.WriteRune('\b')
			case 'f':
				b.WriteRune('\f')
			case 'n':
				b.WriteRune('\n')
			case 't':
				b.WriteRune('\t')
			case 'v':
				b.WriteRune('\v')
			case 'r':
				b.WriteRune('\r')
			case '\\':
				b.WriteRune('\\')
			case '"':
				b.WriteRune('"')
			case '\'':
				b.WriteRune('\'')
			case '0':
				b.WriteRune('\x00')
			case 'x':
				if i+2 < len(s) {
					num, err := strconv.ParseInt(s[i+1:i+3], 16, 8)
					if err == nil {
						b.WriteByte(byte(num))
						i += 2
					} else {
						b.WriteRune('\\')
						b.WriteByte(c)
					}
				} else {
					b.WriteRune('\\')
					b.WriteByte(c)
				}
			case 'o':
				if i+3 < len(s) {
					num, err := strconv.ParseInt(s[i+1:i+4], 8, 8)
					if err == nil {
						b.WriteByte(byte(num))
						i += 3
					} else {
						b.WriteRune('\\')
						b.WriteByte(c)
					}
				} else {
					b.WriteRune('\\')
					b.WriteByte(c)
				}
			default:
				b.WriteRune('\\')
				b.WriteByte(c)
			}
			escaped = false
		} else if c == '\\' {
			escaped = true
		} else {
			b.WriteByte(c)
		}
	}

	return b.String()
}
