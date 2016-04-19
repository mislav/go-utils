package cli

import (
	"os"
)

// Cmd is the environment for each command of the CLI. It provides Stdout, Stderr
// as well as the given flags
type Cmd struct {
	Args   *Args
	Flags  *Flags
	Stdout *ColoredWriter
	Stderr *ColoredWriter
	Env    map[string]string
}

// Exit should be called at the end of each command of the CLI to exit with the
// correct code
func (c *Cmd) Exit(code int) {
	os.Exit(code)
}

// NewCmd initializes a Cmd environment with the given arguments and flags.
// Stdout and Stderr will be initialized with a ColoredWriter on the correct stream
func NewCmd(args *Args, flags *Flags) *Cmd {
	return &Cmd{
		Flags:  flags,
		Stdout: NewColoredWriter(os.Stdout),
		Stderr: NewColoredWriter(os.Stderr),
		Env:    make(map[string]string),
		Args:   args,
	}
}
