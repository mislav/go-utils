package cli

import (
	"os"
)

// Cmd is the environment for each command of the CLI. It provides Stdout, Stderr
// as well as the given flags and arguments
type Cmd struct {
	Args   *Args
	Flags  *Flags
	Stdout *ColoredWriter
	Stderr *ColoredWriter
	Env    CommandConfig
}

// CommandConfig can store additional fields you want to pass to each command
type CommandConfig interface {
}

// Exit should be called at the end of each command of the CLI to exit with the
// correct code
func (c *Cmd) Exit(code int) {
	os.Exit(code)
}

// NewCmd initializes a Cmd environment with the given arguments and flags.
// Stdout and Stderr will be initialized with a ColoredWriter on the correct stream
func NewCmd(args *Args, flags *Flags) *Cmd {
	stderr := NewColoredWriter(os.Stderr)
	stderr.PushColor("red")
	return &Cmd{
		Flags:  flags,
		Stdout: NewColoredWriter(os.Stdout),
		Stderr: stderr,
		Args:   args,
	}
}
