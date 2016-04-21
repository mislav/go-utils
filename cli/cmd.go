package cli

import (
	"os"
)

// Cmd is the environment for each command of the CLI. It provides Stdout, Stderr
// as well as the given flags and arguments
type Cmd struct {
	Args       *Args
	Parameters *Parameters
	Stdout     *ColoredWriter
	Stderr     *ColoredWriter
	Env        CommandConfig
}

// CommandConfig can store additional fields you want to pass to each command
type CommandConfig interface {
}

// ExitValue is a type which describes the return value of the application
type ExitValue int

const (
	// Success lets the programm return with 0
	Success ExitValue = 0
	// Failure lets the programm return with 1
	Failure ExitValue = 1
)

// NewCmd initializes a Cmd environment with the given arguments and flags.
// Stdout and Stderr will be initialized with a ColoredWriter on the correct stream
func NewCmd(args *Args, parameters *Parameters) *Cmd {
	stderr := NewColoredWriter(os.Stderr)
	stderr.PushColor("red")
	return &Cmd{
		Parameters: parameters,
		Stdout:     NewColoredWriter(os.Stdout),
		Stderr:     stderr,
		Args:       args,
	}
}
