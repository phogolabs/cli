package cli

import (
	"fmt"
	"strings"
)

//go:generate parcello -d ./template

var (
	// HelpCommand represents a help command
	HelpCommand = &Command{
		Name:            "help",
		Aliases:         []string{"h"},
		Usage:           "Shows a list of commands or help for one command",
		ArgsUsage:       "[command]",
		HideHelp:        true,
		SkipFlagParsing: true,
		Action:          help,
	}

	// VersionCommand represents a version command
	VersionCommand = &Command{
		Name:            "version",
		Aliases:         []string{"v"},
		Usage:           "Prints the version",
		ArgsUsage:       "[command]",
		HideHelp:        true,
		Hidden:          true,
		SkipFlagParsing: true,
		Action:          version,
	}
)

// Command is a command for a cli.App.
type Command struct {
	// The name of the command
	Name string
	// A list of aliases for the command
	Aliases []string
	// A short description of the usage of this command
	Usage string
	// Custom text to show on USAGE section of help
	UsageText string
	// A longer explanation of how the command works
	Description string
	// A short description of the arguments of this command
	ArgsUsage string
	// The category the command is part of
	Category string
	// Boolean to hide this command from help or completion
	Hidden bool
	// Full name of command for help, defaults to full command name, including parent commands.
	HelpName string
	// Boolean to hide built-in help command
	HideHelp bool
	// List of child commands
	Commands []*Command
	// Treat all flags as normal arguments if true
	SkipFlagParsing bool
	// List of flags to parse
	Flags []Flag
	// Providers contains a list of all providers
	Providers []Provider
	// An action to execute before any subcommands are run, but after the context is ready
	// If a non-nil error is returned, no subcommands are run
	Before BeforeFunc
	// An action to execute after any subcommands are run, but after the subcommand has finished
	// It is run even if Action() panics
	After AfterFunc
	// An action to execute before provider execution
	BeforeInit BeforeFunc
	// An action to execute after provider execution
	AfterInit AfterFunc
	// The action to execute when no subcommands are specified
	// Expects a `cli.ActionFunc` but will accept the *deprecated* signature of `func(*cli.Context) {}`
	Action ActionFunc
	// Execute this function if a usage error occurs.
	OnUsageError OnUsageErrorFunc
	// Metadata information
	Metadata map[string]interface{}
}

// RunWithContext runs the command
func (cmd *Command) RunWithContext(ctx *Context) error {
	cmd.prepare()

	if err := cmd.provide(ctx); err != nil {
		return cmd.error(ctx, err)
	}

	if err := cmd.fork(ctx); err != ErrCommandNotFound {
		return err
	}

	return cmd.exec(ctx)
}

// Names returns the names including short names and aliases.
func (cmd *Command) Names() []string {
	names := []string{cmd.Name}
	return append(names, cmd.Aliases...)
}

// VisibleFlags returns a slice of the Flags with Hidden=false
func (cmd *Command) VisibleFlags() []Flag {
	flags := []Flag{}

	for _, flag := range cmd.Flags {
		accessor := NewFlagAccessor(flag)

		if accessor.Hidden() {
			continue
		}

		flags = append(flags, accessor)
	}

	return flags
}

// VisibleCommands returns a slice of the Commands with Hidden=false
func (cmd *Command) VisibleCommands() []*Command {
	category := &CommandCategory{
		Commands: cmd.Commands,
	}

	return category.VisibleCommands()
}

// VisibleCategories returns a slice of categories and commands that are
// Hidden=false
func (cmd *Command) VisibleCategories() []*CommandCategory {
	result := []*CommandCategory{}
	categories := map[string]*CommandCategory{}

	for _, command := range cmd.Commands {
		name := command.Category
		category, ok := categories[name]

		if !ok {
			category = &CommandCategory{Name: name}
			result = append(result, category)
			categories[name] = category
		}

		category.Commands = append(category.Commands, command)
	}

	return result
}

func (cmd *Command) provide(ctx *Context) (err error) {
	if cmd.SkipFlagParsing {
		return nil
	}

	if cmd.AfterInit != nil {
		defer func() {
			if afterErr := cmd.AfterInit(ctx); afterErr != nil {
				err = AppendError(err, afterErr)
			}
		}()
	}

	if cmd.BeforeInit != nil {
		if beforeErr := cmd.BeforeInit(ctx); beforeErr != nil {
			err = AppendError(err, beforeErr)
		}
	}

	if err != nil {
		return err
	}

	for _, provider := range cmd.Providers {
		if err = provider.Provide(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Command) validate(ctx *Context) error {
	for _, flag := range cmd.Flags {
		if err := flag.Validate(); err != nil {
			return cmd.error(ctx, err)
		}
	}

	return nil
}

func (cmd *Command) prepare() {
	cmd.providers()
	cmd.flags()
	cmd.commands()
}

func (cmd *Command) providers() {
	providers := []Provider{
		&FileProvider{},
		&EnvProvider{},
		&FlagProvider{},
	}

	cmd.Providers = append(providers, cmd.Providers...)
}

func (cmd *Command) commands() {
	if cmd.HelpName == "" {
		cmd.HelpName = cmd.Name
	}

	if !cmd.HideHelp {
		cmd.Commands = append(cmd.Commands, HelpCommand)
	}

	for _, command := range cmd.Commands {
		if command.HelpName == "" {
			command.HelpName = fmt.Sprintf("%s %s", cmd.HelpName, command.Name)
		}
	}
}

func (cmd *Command) flags() {
	if !cmd.HideHelp {
		help := &BoolFlag{
			Name:  "help, h",
			Usage: "shows help",
		}

		cmd.Flags = append(cmd.Flags, help)
	}

	if cmd.Metadata == nil {
		cmd.Metadata = make(map[string]interface{})
	}

	cmd.Metadata["VisibleFlags"] = cmd.VisibleFlags()
}

func (cmd *Command) fork(ctx *Context) error {
	var (
		child *Command
		args  []string
	)

	switch {
	case ctx.Bool("help"):
		child = cmd.find("help")
	case ctx.Bool("version"):
		child = cmd.find("version")
	case len(ctx.Args) > 0:
		child, args = cmd.next(ctx.Args)
	case cmd.Action == nil:
		child = cmd.find("help")
	}

	if child == nil {
		return ErrCommandNotFound
	}

	switch {
	case child.Name == "help":
		break
	case child.Name == "version":
		break
	default:
		if err := cmd.validate(ctx); err != nil {
			return err
		}
	}

	ctx = &Context{
		Parent:    ctx,
		Metadata:  ctx.Metadata,
		Manifest:  ctx.Manifest,
		Writer:    ctx.Writer,
		ErrWriter: ctx.ErrWriter,
		Command:   child,
		Args:      args,
	}

	return child.RunWithContext(ctx)
}

func (cmd *Command) next(args []string) (*Command, []string) {
	child := cmd.find(args[0])

	if child == nil {
		child = cmd.find("help")
		args = args[:1]
	} else {
		args = args[1:]
	}

	return child, args
}

func (cmd *Command) find(name string) *Command {
	for _, child := range cmd.Commands {
		for _, alias := range child.Names() {
			if strings.EqualFold(alias, name) {
				return child
			}
		}
	}

	return nil
}

func (cmd *Command) exec(ctx *Context) (err error) {
	if cmd.After != nil {
		defer func() {
			if afterErr := cmd.After(ctx); afterErr != nil {
				err = AppendError(err, afterErr)
			}
		}()
	}

	if cmd.Before != nil {
		if beforeErr := cmd.Before(ctx); beforeErr != nil {
			err = AppendError(err, beforeErr)
		}
	}

	if err != nil {
		return err
	}

	if err = cmd.validate(ctx); err != nil {
		return err
	}

	return cmd.Action(ctx)
}

func (cmd *Command) error(ctx *Context, err error) error {
	if cmd.OnUsageError != nil {
		err = cmd.OnUsageError(ctx, err)
	}

	if err == nil {
		return err
	}

	fmt.Fprintln(ctx.Writer, "Incorrect Usage:", err.Error())
	fmt.Fprintln(ctx.Writer)

	ctx.Args = []string{"help"}

	cmd.fork(ctx)

	return err
}
