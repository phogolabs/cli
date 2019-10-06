package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/phogolabs/parcello"
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

	var (
		help, _         = parcello.Open(man)
		content, _      = ioutil.ReadAll(help)
		writer          = tabwriter.NewWriter(ctx.Writer, 1, 8, 2, ' ', 0)
		templateFuncMap = template.FuncMap{
			"join": strings.Join,
		}
	)

	tmpl := template.New("help").Funcs(templateFuncMap)
	tmpl = template.Must(tmpl.Parse(string(content)))

	if err := tmpl.Execute(writer, cmd); err != nil {
		return err
	}

	return writer.Flush()
}

func version(ctx *Context) error {
	for ctx.Parent != nil {
		ctx = ctx.Parent
	}

	tmpl := template.New("version")
	tmpl = template.Must(tmpl.Parse("{{ .Name }} version {{ .Metadata.Version }}\n"))

	if err := tmpl.Execute(ctx.Writer, ctx.Command); err != nil {
		return err
	}

	return nil
}
