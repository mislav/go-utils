package cli

import (
	"path"
	"regexp"
	"strings"
)

type Args struct {
	programName string
	argv        []string
}

func NewArgs(argv []string) *Args {
	return &Args{
		programName: argv[0],
		argv:        argv[1:],
	}
}

func (a *Args) SubcommandArgs(cmdName string) *Args {
	var newArgv []string
	if len(a.argv) > 1 {
		newArgv = a.argv[1:]
	} else {
		newArgv = []string{}
	}

	newProgramName := a.programName
	if cmdName != "" {
		newProgramName += " " + cmdName
	}

	return &Args{
		programName: newProgramName,
		argv:        newArgv,
	}
}

func (a *Args) ProgramName() string {
	return path.Base(a.programName)
}

func (a *Args) Length() int {
	return len(a.argv)
}

func (a *Args) Peek(index int) string {
	if len(a.argv) > index {
		return a.argv[index]
	}
	return ""
}

func (a *Args) Shift() (string, *Args) {
	if !strings.HasPrefix(a.argv[0], "-") {
		return a.argv[0], a.SubcommandArgs("")
	}
	return "", a
}

func (a *Args) Slice(start int) []string {
	return a.argv[start:]
}

func (a *Args) String() string {
	return strings.Join(a.argv, " ")
}

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
