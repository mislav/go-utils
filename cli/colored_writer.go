package cli

import (
	"fmt"
	"io"
	"os"

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
