package cli

import "sync"

var appSingleton *App
var once sync.Once

// App saves registered subcommands, flags and a bit more
type App struct {
	Name               string
	DefaultCommandName string
	commands           map[string]Command
	Fallback           func(c *Cmd, cmdName string)
	Before             func(c *Cmd)
	flags              map[string]*Flag
}

// Command is the struct for one subcommand of the app including help text and flags
type Command struct {
	Name     string
	Help     string
	Function func(*Cmd)
	flags    map[string]*Flag
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
func (a *App) Flags() map[string]*Flag {
	return a.flags
}

// RegisterCommand registers a Command
func (a *App) RegisterCommand(c Command) {
	if a.commands == nil {
		a.commands = make(map[string]Command)
	}
	a.commands[c.Name] = c
}

// RegisterFlag registers a Flag
func (a *App) RegisterFlag(f Flag) {
	if a.flags == nil {
		a.flags = make(map[string]*Flag)
	}
	a.flags[f.Long] = &f
}

// RegisterFlag registers a Flag
func (c *Command) RegisterFlag(f Flag) {
	if c.flags == nil {
		c.flags = make(map[string]*Flag)
	}
	c.flags[f.Long] = &f
}

// Run executs the App with the given arguments
func (a *App) Run(arguments []string) {
	args := NewArgs(arguments)
	cmdName := args.Peek(0)
	if cmdName == "" {
		cmdName = a.DefaultCommandName
	}
	command := a.commands[cmdName]
	cmdFunc := command.Function
	flags := &Flags{}
	for _, flag := range a.flags {
		flag, args = args.Extract(*flag)
		flags.AddFlag(flag)
	}
	for _, flag := range command.flags {
		flag, args = args.Extract(*flag)
		flags.AddFlag(flag)
	}
	cmd := NewCmd(args, flags)
	if a.Before != nil {
		a.Before(cmd)
	}
	if cmdFunc != nil {
		cmdFunc(cmd)
	} else {
		a.Fallback(cmd, cmdName)
	}
}
