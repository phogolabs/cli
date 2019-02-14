package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/phogolabs/parcello"
)

// TemplateFuncMap exposes a map of function used in the templates
var TemplateFuncMap = template.FuncMap{
	"join": strings.Join,
}

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

	cmd.restore(ctx)

	help, _ := parcello.Open(man)
	content, _ := ioutil.ReadAll(help)
	writer := tabwriter.NewWriter(ctx.Writer, 1, 8, 2, ' ', 0)

	tmpl := template.New("help").Funcs(TemplateFuncMap)
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

	fmt.Fprintf(ctx.Writer, "%v version %v\n", ctx.Command.Name, ctx.Command.Metadata["Version"])
	return nil
}
