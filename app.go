package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path"
	"time"
)

// App is the main structure of a cli application. It is recommended that
// an app be created with the cli.NewApp() function
type App struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string
	// A short description of the usage of this command
	Usage string
	// Custom text to show on USAGE section of help
	UsageText string
	// A longer explanation of how the command works
	Description string
	// A short description of the arguments of this command
	ArgsUsage string
	// Compilation date
	Compiled time.Time
	// List of all authors who contributed
	Authors []*Author
	// Copyright of the binary if any
	Copyright string
	// Full name of command for help, defaults to full command name, including parent commands.
	HelpName string
	// Boolean to hide built-in help command
	HideHelp bool
	// Version of the program
	Version string
	// Boolean to hide built-in version flag and the VERSION section of help
	HideVersion bool
	// Signals are the signals that we want to handle
	Signals []os.Signal
	// List of commands to execute
	Commands []*Command
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
	// OnSignal occurs on system signal
	OnSignal SignalFunc
	// Execute this function if a usage error occurs.
	OnUsageError UsageErrorFunc
	// OnCommandNotFound is executed if the proper command cannot be found
	OnCommandNotFound CommandNotFoundFunc
	// Execute this function to handle ExitErrors. If not provided, HandleExitCoder is provided to
	// function as a default, so this is optional.
	OnExitError ExitErrorHandlerFunc
	// Exit is the function used when the app exits. If not set defaults to os.Exit.
	Exit ExitFunc
	// Writer writer to write output to
	Writer io.Writer
	// ErrWriter writes error output
	ErrWriter io.Writer
}

// Run is the entry point to the cli app. Parses the arguments slice and routes
// to the proper flag/args combination
func (app *App) Run(args []string) {
	args = app.prepare(args)

	cmd := &Command{
		Name:              app.Name,
		Usage:             app.Usage,
		UsageText:         app.UsageText,
		HideHelp:          app.HideHelp,
		HelpName:          app.HelpName,
		Commands:          app.Commands,
		Description:       app.Description,
		ArgsUsage:         app.ArgsUsage,
		Flags:             app.Flags,
		Before:            app.Before,
		After:             app.After,
		BeforeInit:        app.BeforeInit,
		AfterInit:         app.AfterInit,
		Action:            app.Action,
		Providers:         app.Providers,
		OnUsageError:      app.OnUsageError,
		OnCommandNotFound: app.OnCommandNotFound,
		Metadata: Map{
			"HideVersion": app.HideVersion,
			"Version":     app.Version,
			"Authors":     app.Authors,
			"Copyright":   app.Copyright,
		},
	}

	ctx := &Context{
		Command:   cmd,
		Args:      args[1:],
		Writer:    app.Writer,
		ErrWriter: app.ErrWriter,
		Metadata:  make(map[string]interface{}),
	}

	app.notify(ctx)
	app.error(cmd.RunWithContext(ctx))
}

func (app *App) notify(ctx *Context) {
	if len(app.Signals) == 0 {
		return
	}

	if app.OnSignal == nil {
		return
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, app.Signals...)

	go func() {
		signal := <-ch

		err := app.OnSignal(ctx, signal)
		app.error(err)
	}()
}

func (app *App) prepare(args []string) []string {
	app.flags()
	app.commands()
	return app.app(args)
}

func (app *App) app(args []string) []string {
	if len(args) == 0 {
		args = []string{"unknown"}
	}

	if app.Name == "" {
		app.Name = path.Base(args[0])
	}

	if app.Compiled.IsZero() {
		info, err := os.Stat(args[0])

		if err != nil {
			app.Compiled = time.Now()
		} else {
			app.Compiled = info.ModTime()
		}
	}

	if app.Writer == nil {
		app.Writer = os.Stdout
	}

	if app.ErrWriter == nil {
		app.ErrWriter = os.Stderr
	}

	if app.Exit == nil {
		app.Exit = os.Exit
	}

	return args
}

func (app *App) flags() {
	if !app.HideVersion {
		version := &BoolFlag{
			Name:  "version, v",
			Usage: "prints the version",
		}

		app.Flags = append(app.Flags, version)
	}
}

func (app *App) commands() {
	if !app.HideVersion {
		version := NewVersionCommand()
		app.Commands = append(app.Commands, version)
	}
}

func (app *App) error(err error) {
	if err == nil {
		return
	}

	if app.OnExitError != nil {
		err = app.OnExitError(err)
	}

	fmt.Fprintln(app.ErrWriter, err)

	var errx *ExitError

	if !errors.As(err, &errx) {
		errx = WrapError(err)
	}

	app.Exit(errx.Code())
}

// Author represents someone who has contributed to a cli project.
type Author struct {
	// Name of the author
	Name string
	// Email of the author
	Email string
}

// String makes Author comply to the Stringer interface, to allow an easy print in the templating process
func (author *Author) String() string {
	value := ""

	if author.Email != "" {
		value = fmt.Sprintf(" <%s>", author.Email)
	}

	return fmt.Sprintf("%v%v", author.Name, value)
}

// Copyright creates a copyright message
func Copyright(name string) string {
	return fmt.Sprintf("%s (C) %d", name, time.Now().Year())
}
