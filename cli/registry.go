package cli

var commands = make(map[string]func(*Cmd))
var help = make(map[string]string)

func Lookup(cmdName string) func(*Cmd) {
	return commands[cmdName]
}

func LookupHelp(cmdName string) string {
	return help[cmdName]
}

func Register(cmdName, helpString string, fn func(*Cmd)) {
	commands[cmdName] = fn
	help[cmdName] = helpString
}

func CommandNames() []string {
	names := make([]string, len(commands))
	i := 0
	for name, _ := range commands {
		names[i] = name
		i += 1
	}
	return names
}
