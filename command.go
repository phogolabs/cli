package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/phogolabs/log"
)

//go:generate parcello -d ./template

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
	// Metadata information
	Metadata map[string]interface{}
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
	After AfterFunc
	// An action to execute before provider execution
	BeforeInit BeforeFunc
	// An action to execute after provider execution
	AfterInit AfterFunc
	// The action to execute when no subcommands are specified
	// Expects a cli.ActionFunc
	Action ActionFunc
	// Strategy enables comman retry logic
	Strategy BackOffStrategy
	// Execute this function if a usage error occurs.
	OnUsageError UsageErrorFunc
	// OnCommandNotFound is executed if the proper command cannot be found
	OnCommandNotFound CommandNotFoundFunc
}

// NewHelpCommand creates a new help command
func NewHelpCommand() *Command {
	help := &Command{
		Name:            "help",
		Aliases:         []string{"h"},
		Usage:           "Shows a list of commands or help for one command",
		ArgsUsage:       "[command]",
		HideHelp:        true,
		SkipFlagParsing: true,
		Action:          help,
	}

	return help
}

// NewVersionCommand creates a version command
func NewVersionCommand() *Command {
	return &Command{
		Name:            "version",
		Aliases:         []string{"v"},
		Usage:           "Prints the version",
		ArgsUsage:       "[command]",
		HideHelp:        true,
		Hidden:          true,
		SkipFlagParsing: true,
		Action:          version,
	}
}

// RunWithContext runs the command
func (cmd *Command) RunWithContext(ctx *Context) error {
	cmd.prepare()

	if err := cmd.provide(ctx); err != nil {
		return cmd.error(ctx, err)
	}

	err := cmd.fork(ctx)

	if errx, ok := err.(ExitCoder); ok {
		if errx.Code() == ExitCodeNotFoundCommand {
			return cmd.exec(cmd.Action, ctx)
		}
	}

	return err
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

func (cmd *Command) provide(ctx *Context) (errx error) {
	if cmd.SkipFlagParsing {
		return nil
	}

	var errs ExitErrorCollector

	defer func() {
		if len(errs) > 0 {
			errx = errs
		}
	}()

	if cmd.AfterInit != nil {
		defer func() {
			if afterErr := cmd.retry(cmd.AfterInit, ctx); afterErr != nil {
				errs.Wrap(afterErr)
			}
		}()
	}

	if cmd.BeforeInit != nil {
		if beforeErr := cmd.retry(cmd.BeforeInit, ctx); beforeErr != nil {
			errs.Wrap(beforeErr)
			return
		}
	}

	for _, provider := range cmd.Providers {
		if err := cmd.retry(provider.Provide, ctx); err != nil {
			errs.Wrap(err)
			return
		}
	}

	return
}

func (cmd *Command) validate(ctx *Context) error {
	for _, flag := range cmd.Flags {
		accessor := NewFlagAccessor(flag)
		if err := accessor.Validate(ctx); err != nil {
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
		help := NewHelpCommand()
		cmd.Commands = append(cmd.Commands, help)
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
		name  string
		args  []string
	)

	switch {
	case ctx.Bool("help"):
		child = cmd.find("help")
	case ctx.Bool("version"):
		child = cmd.find("version")
	case len(ctx.Args) > 0:
		name = ctx.Args[0]
		child, args = cmd.next(ctx.Args)
	case cmd.Action == nil:
		child = cmd.find("help")
	}

	if child == nil {
		if cmd.OnCommandNotFound != nil {
			cmd.OnCommandNotFound(ctx, name)
		}

		return NotFoundCommandError(name)
	}

	ctx = &Context{
		Parent:    ctx,
		Metadata:  ctx.Metadata,
		Writer:    ctx.Writer,
		ErrWriter: ctx.ErrWriter,
		Command:   child,
		Args:      args,
	}

	if child.Name != "help" && child.Name != "version" {
		return cmd.exec(child.RunWithContext, ctx)
	}

	return child.RunWithContext(ctx)
}

func (cmd *Command) next(args []string) (*Command, []string) {
	child := cmd.find(args[0])

	if child == nil {
		if !cmd.has() {
			return nil, []string{}
		}

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

func (cmd *Command) exec(action ActionFunc, ctx *Context) (errx error) {
	var (
		errs     ExitErrorCollector
		eventCtx *Context = ctx
	)

	defer func() {
		if !errs.IsEmpty() {
			errx = errs
		}
	}()

	if ctx.Command != cmd {
		eventCtx = ctx.Parent
	}

	if cmd.After != nil {
		defer func() {
			if afterErr := cmd.retry(cmd.After, eventCtx); afterErr != nil {
				errs.Wrap(afterErr)
			}
		}()
	}

	if cmd.Before != nil {
		if beforeErr := cmd.retry(cmd.Before, eventCtx); beforeErr != nil {
			errs.Wrap(beforeErr)
			return
		}
	}

	if err := cmd.validate(eventCtx); err != nil {
		errs.Wrap(err)
		return
	}

	if err := cmd.retry(action, ctx); err != nil {
		errs.Wrap(err)
		return
	}

	return nil
}

func (cmd *Command) retry(fn func(*Context) error, ctx *Context) error {
	if cmd.Strategy == nil {
		return fn(ctx)
	}

	logger := log.WithField("command", cmd.Name)

	tryFunc := func() error {
		return fn(ctx)
	}

	notify := func(err error, t time.Duration) {
		logger.WithError(err).Warnf("executing the command not successful. retry in %v", t)
	}

	if err := backoff.RetryNotify(tryFunc, cmd.Strategy, notify); err != nil {
		logger.WithError(err).Fatal("executing the command failed")
		return err
	}

	return nil
}

func (cmd *Command) has() bool {
	count := len(cmd.Commands)
	for _, sub := range cmd.Commands {
		switch {
		case sub.Name == "help":
			count--
		case sub.Name == "version":
			count--
		}
	}

	return count > 0
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

// CommandsByName is a slice of Command
type CommandsByName []*Command

// Len returns the length of the slice
func (c CommandsByName) Len() int {
	return len(c)
}

// Less returns true if item at index i < item at index j
func (c CommandsByName) Less(i, j int) bool {
	return lexicographicLess(c[i].Name, c[j].Name)
}

// Swap swaps two items
func (c CommandsByName) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
