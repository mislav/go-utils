package cli

var commands = make(map[string]func(*Cmd))

func Lookup(cmdName string) func(*Cmd) {
	return commands[cmdName]
}

func Register(cmdName string, fn func(*Cmd)) {
	commands[cmdName] = fn
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
