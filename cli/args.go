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
	} else {
		return ""
	}
}

func (a *Args) Shift() (string, *Args) {
	if !strings.HasPrefix(a.argv[0], "-") {
		return a.argv[0], a.SubcommandArgs("")
	} else {
		return "", a
	}
}

func (a *Args) String() string {
	return strings.Join(a.argv, " ")
}

type Flag struct {
	Short    string
	Long     string
	values   []string
	provided bool
}

func (f *Flag) AddValue(v string) {
	f.values = append(f.values, v)
}

func (f *Flag) IsProvided() bool {
	return f.provided
}

func (f *Flag) String() string {
	num := len(f.values)
	if num > 0 {
		return f.values[num-1]
	} else {
		return ""
	}
}

func (f *Flag) Strings() []string {
	return f.values
}

func (f *Flag) Bool() bool {
	val := strings.ToLower(f.String())
	return val == "true" || val == "t" || val == "1"
}

func (a *Args) ExtractFlag(short, long string, ftype interface{}) (*Flag, *Args) {
	explodedArgv := []string{}
	multipleShortRE := regexp.MustCompile(`^-[a-zA-Z]{2,}$`)
	flag := &Flag{
		Short: short,
		Long:  long,
	}

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
	_, isBool := ftype.(bool)
	grabNextValue := false

	for _, arg := range explodedArgv {
		if grabNextValue {
			flag.AddValue(arg)
			grabNextValue = false
			continue
		}

		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg, "=", 2)
			if flag.Long != "" && parts[0] == flag.Long {
				if len(parts) > 1 {
					flag.AddValue(parts[1])
				} else if isBool {
					flag.AddValue("true")
				} else {
					grabNextValue = true
				}
				flag.provided = true
				continue
			}
		} else if flag.Short != "" && arg == flag.Short {
			if isBool {
				flag.AddValue("true")
			} else {
				grabNextValue = true
			}
			flag.provided = true
			continue
		}

		newArgv = append(newArgv, arg)
	}

	args := &Args{
		programName: a.programName,
		argv:        newArgv,
	}

	return flag, args
}
