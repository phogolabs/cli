package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/phogolabs/cli/template"
)

func help(ctx *Context) error {
	var (
		man  string
		name string
		cmd  *Command
	)

	switch {
	case len(ctx.Args) > 0:
		name = ctx.Args[0]
		cmd = ctx.Parent.Command.find(name)
		man = "help.cmd.tpl"
	case ctx.Parent != nil:
		cmd = ctx.Parent.Command
		ctx = ctx.Parent
		man = "help.sub.tpl"
	default:
		cmd = ctx.Command
		man = "help.app.tpl"
	}

	if cmd == nil {
		fmt.Fprintf(ctx.Writer, "No help topic for '%s'", name)
		fmt.Fprintln(ctx.Writer)
		return nil
	}

	writer := tabwriter.NewWriter(ctx.Writer, 1, 8, 2, ' ', 0)

	content, err := template.Open(man)
	if err != nil {
		return err
	}

	if err := content.Execute(writer, cmd); err != nil {
		return err
	}

	return writer.Flush()
}

func version(ctx *Context) error {
	for ctx.Parent != nil {
		ctx = ctx.Parent
	}

	content, err := template.Open("version.app.tpl")
	if err != nil {
		return err
	}

	if err := content.Execute(ctx.Writer, ctx.Command); err != nil {
		return err
	}

	return nil
}
