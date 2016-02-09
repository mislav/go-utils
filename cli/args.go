package cli

import (
	"path"
	"strings"
)

type Args struct {
	programName string
	argv []string
}

func NewArgs(argv []string) *Args {
	return &Args{
		programName: argv[0],
		argv: argv[1:],
	}
}

func (a *Args) SubcommandArgs() *Args {
	return &Args{
		programName: a.programName + " " + a.argv[0],
		argv: a.argv[1:],
	}
}

func (a *Args) ProgramName() string {
	return path.Base(a.programName)
}

func (a *Args) At(n int) string {
	if len(a.argv) > n {
		return a.argv[n]
	} else {
		return ""
	}
}

func (a *Args) Word(n int) string {
	for {
		arg := a.At(n)
		if arg != "" && strings.HasPrefix(arg, "-") {
			n += 1
			continue
		}
		return arg
	}
}

func (a *Args) HasFlag(flag string) bool {
	for _, arg := range a.argv {
		if arg == flag {
			return true
		}
	}
	return false
}
