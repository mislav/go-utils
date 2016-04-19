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
	stderr := NewColoredWriter(os.Stderr)
	stderr.PushColor("red")
	return &Cmd{
		Args:   args,
		Stdout: NewColoredWriter(os.Stdout),
		Stderr: stderr,
	}
}
