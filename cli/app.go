package cli

import "sync"

var appSingleton *App
var once sync.Once

// App saves registered commands, flags and a bit more
type App struct {
	// DefaultCommandName: a command with this name will be called if no other name was provided
	DefaultCommandName string
	Version            string
	commands           map[string]Command
	// Fallback function will be called if no command with the given name was found
	Fallback func(c *Cmd, cmdName string) ExitValue
	// Before function will be called after cmd is initialize and can adjust it
	// before it will be passed to the correct command
	Before func(c *Cmd, cmdName string)
	flags  map[string]Flag
}

// Command is the struct for one command of the app including help text and flags
type Command struct {
	Name string
	// Info is a short string shown behind the command in the help
	Info string
	// Help text will be shown after the info text if the user requests help for this command
	Help string
	// Parameter string which will be displayed in the usage section of the help
	Parameter string
	Function  func(*Cmd) ExitValue
	flags     map[string]Flag
	commands  map[string]Command
}

// AppInstance returns the singleton instance of App
func AppInstance() *App {
	once.Do(func() {
		appSingleton = &App{}
	})
	return appSingleton
}

// Commands returns all registered Commands
func (a *App) Commands() map[string]Command {
	return a.commands
}

// Flags returns all registered application wide Flags
func (a *App) Flags() map[string]Flag {
	return a.flags
}

// RegisterCommand registers Commands
func (a *App) RegisterCommand(c ...Command) {
	if a.commands == nil {
		a.commands = make(map[string]Command)
	}
	for _, command := range c {
		a.commands[command.Name] = command
	}
}

// RegisterFlag registers Flags
func (a *App) RegisterFlag(f ...Flag) {
	if a.flags == nil {
		a.flags = make(map[string]Flag)
	}
	for _, flag := range f {
		a.flags[flag.Long] = flag
	}
}

// Commands returns all registered Commands
func (c *Command) Commands() map[string]Command {
	return c.commands
}

// RegisterCommand registers Commands
func (c *Command) RegisterCommand(cmds ...Command) {
	if c.commands == nil {
		c.commands = make(map[string]Command)
	}
	for _, command := range cmds {
		c.commands[command.Name] = command
	}
}

// RegisterFlag registers Flags
func (c *Command) RegisterFlag(f ...Flag) {
	if c.flags == nil {
		c.flags = make(map[string]Flag)
	}
	for _, flag := range f {
		c.flags[flag.Long] = flag
	}
}

// Flags returns all registered Flags for the Command
func (c *Command) Flags() map[string]Flag {
	return c.flags
}

// Run executs the App with the given arguments
func (a *App) Run(arguments []string) ExitValue {
	args := NewArgs(arguments)
	cmdName := args.Peek(0)
	if cmdName == "" {
		cmdName = a.DefaultCommandName
	}
	command := a.commands[cmdName]
	parameters := &Parameters{}
	var parameter *Parameter
	for _, flag := range a.flags {
		parameter, args = args.Extract(flag)
		parameters.AddParameter(parameter)
	}
	args, cmdFunc := a.command(command, args, parameters)
	cmd := NewCmd(args, parameters)
	if a.Before != nil {
		a.Before(cmd, cmdName)
	}
	if cmdFunc != nil {
		return cmdFunc(cmd)
	}
	return a.Fallback(cmd, cmdName)
}

func (a *App) command(command Command, args *Args, parameters *Parameters) (*Args, func(*Cmd) ExitValue) {
	args = args.SubcommandArgs(command.Name)
	var parameter *Parameter
	for _, flag := range command.flags {
		parameter, args = args.Extract(flag)
		parameters.AddParameter(parameter)
	}
	if command.commands != nil {
		subCommandName := args.Peek(0)
		if subCommand, ok := command.commands[subCommandName]; ok {
			return a.command(subCommand, args, parameters)
		}
	}
	return args, command.Function
}
