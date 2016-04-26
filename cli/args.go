package cli

import (
	"path"
	"regexp"
	"strings"
)

// Args holds the arguments the cli recived
type Args struct {
	programName string
	CommandName string
	argv        []string
}

// NewArgs constructs Args form a given cmd argv
func NewArgs(argv []string) *Args {
	return &Args{
		programName: argv[0],
		CommandName: "",
		argv:        argv[1:],
	}
}

// SubcommandArgs returns the arguments for the given subcommand.
// Notice: If the given name is not the called subcommand it has no arguments
func (a *Args) SubcommandArgs(cmdName string) *Args {
	var newArgv []string
	if len(a.argv) > 1 {
		if a.Peek(0) != cmdName {
			newArgv = []string{}
			cmdName = ""
		}
		newArgv = a.argv[1:]
	} else {
		newArgv = []string{}
	}

	newCommandName := a.CommandName
	if newCommandName == "" {
		newCommandName = cmdName
	} else if cmdName != "" {
		newCommandName += " " + cmdName
	}

	return &Args{
		programName: a.programName,
		CommandName: newCommandName,
		argv:        newArgv,
	}
}

// ProgramName returns the name of the executable
func (a *Args) ProgramName() string {
	return path.Base(a.programName)
}

// Length returns the number of arguments
func (a *Args) Length() int {
	return len(a.argv)
}

// Peek returns the argument at the given index
func (a *Args) Peek(index int) string {
	if len(a.argv) > index {
		return a.argv[index]
	}
	return ""
}

// Shift extracts the first argument and returns the value and the remaining Args
func (a *Args) Shift() (string, *Args) {
	return a.argv[0], a.SubcommandArgs("")
}

//Slice returns arguments from the given start point
func (a *Args) Slice(start int) []string {
	return a.argv[start:]
}

// String returns all arguments seperated as
func (a *Args) String() string {
	return strings.Join(a.argv, " ")
}

// Extract extracts a flag from the arguments
func (a *Args) Extract(flag Flag) (*Parameter, *Args) {
	parameter := Parameter{Flag: flag}
	explodedArgv := []string{}
	multipleShortRE := regexp.MustCompile(`^-[a-zA-Z]{2,}$`)

	for _, arg := range a.argv {
		if multipleShortRE.MatchString(arg) {
			shorts := strings.TrimPrefix(arg, "-")
			for _, f := range strings.Split(shorts, "") {
				explodedArgv = append(explodedArgv, "-"+f)
			}
		} else {
			explodedArgv = append(explodedArgv, arg)
		}
	}

	newArgv := []string{}
	_, isBool := parameter.Ftype.(bool)
	grabNextValue := false

	for _, arg := range explodedArgv {
		if grabNextValue {
			parameter.AddValue(arg)
			grabNextValue = false
			continue
		}

		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg, "=", 2)
			if parameter.Long != "" && parts[0] == parameter.Long {
				if len(parts) > 1 {
					parameter.AddValue(parts[1])
				} else if isBool {
					parameter.AddValue("true")
				} else {
					grabNextValue = true
				}
				parameter.provided = true
				continue
			}
		} else if parameter.Short != "" && arg == parameter.Short {
			if isBool {
				parameter.AddValue("true")
			} else {
				grabNextValue = true
			}
			parameter.provided = true
			continue
		}

		newArgv = append(newArgv, arg)
	}

	args := &Args{
		programName: a.programName,
		argv:        newArgv,
	}

	return &parameter, args
}
