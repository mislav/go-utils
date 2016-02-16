package cli

import (
	"os"
)

type Cmd struct {
	Args   *Args
	Stdout *ColoredWriter
	Stderr *ColoredWriter
}

func (c *Cmd) Exit(code int) {
	os.Exit(code)
}

func NewCmd(args *Args) *Cmd {
	return &Cmd{
		Args:   args,
		Stdout: NewColoredWriter(os.Stdout),
		Stderr: NewColoredWriter(os.Stderr),
	}
}
