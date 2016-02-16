package cli

import (
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

var colors = map[string]string{
	"black":   "30",
	"red":     "31",
	"green":   "32",
	"yellow":  "33",
	"blue":    "34",
	"magenta": "35",
	"cyan":    "36",
	"white":   "37",

	"bold":        "1",
	"boldblack":   "30;1",
	"boldred":     "31;1",
	"boldgreen":   "32;1",
	"boldyellow":  "33;1",
	"boldblue":    "34;1",
	"boldmagenta": "35;1",
	"boldcyan":    "36;1",
	"boldwhite":   "37;1",
}

type ColoredWriter struct {
	colorize   bool
	colorStack []string
	out        io.Writer
}

func NewColoredWriter(file *os.File) *ColoredWriter {
	return &ColoredWriter{
		out:      colorable.NewColorable(file),
		colorize: isatty.IsTerminal(file.Fd()),
	}
}

func (c *ColoredWriter) AlwaysColorize() {
	c.colorize = true
}

func (c *ColoredWriter) PushColor(color string) {
	if c.colorize {
		c.colorStack = append(c.colorStack, color)
	}
}

func (c *ColoredWriter) PopColor() {
	if c.colorize {
		c.colorStack = c.colorStack[0 : len(c.colorStack)-1]
	}
}

func (c *ColoredWriter) CurrentColor() string {
	if len(c.colorStack) > 0 {
		return c.colorStack[len(c.colorStack)-1]
	} else {
		return ""
	}
}

func (c *ColoredWriter) Write(p []byte) (int, error) {
	color := c.CurrentColor()

	if color != "" {
		_, err := fmt.Fprintf(c.out, "\033[%sm", colors[color])
		if err != nil {
			return 0, err
		}
	}

	n, err := c.out.Write(p)

	if err == nil && color != "" {
		_, err = c.out.Write([]byte("\033[0m"))
	}

	return n, err
}

func (c *ColoredWriter) Cprintln(color string, a ...interface{}) (n int, err error) {
	var n2 int
	n, err = c.Cprint(color, a...)
	if err == nil {
		n2, err = fmt.Fprint(c.out, "\n")
		n += n2
	}
	return
}

func (c *ColoredWriter) Cprint(color string, a ...interface{}) (n int, err error) {
	if c.colorize {
		_, err = fmt.Fprintf(c.out, "\033[%sm", colors[color])
		if err != nil {
			return 0, err
		}
	}
	n, err = fmt.Fprint(c.out, a...)
	if c.colorize {
		_, err = c.out.Write([]byte("\033[0m"))
	}
	return
}

var cprintfRe = regexp.MustCompile(`%C\((\w+)\)`)

func (c *ColoredWriter) Cprintf(format string, a ...interface{}) (int, error) {
	if matches := cprintfRe.FindAllStringSubmatch(format, -1); matches != nil {
		i := 0
		format = cprintfRe.ReplaceAllStringFunc(format, func(string) string {
			if !c.colorize {
				return ""
			}
			color := matches[i][1]
			i += 1
			var colorCode string
			if color == "reset" {
				if defaultColor := c.CurrentColor(); defaultColor == "" {
					colorCode = "0"
				} else {
					colorCode = colors[defaultColor]
				}
			} else {
				colorCode = colors[color]
			}
			return fmt.Sprintf("\033[%sm", colorCode)
		})
	}

	return c.Printf(format, a...)
}

func (c *ColoredWriter) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(c, a...)
}

func (c *ColoredWriter) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(c, a...)
}

func (c *ColoredWriter) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(c, format, a...)
}
