NAME:
   {{.Name}}{{if .Usage}} - {{.Usage}}{{end}}
USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Metadata.Version}}{{if not .Metadata.HideVersion}}
VERSION:
   {{.Metadata.Version}}{{end}}{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if len .Metadata.Authors}}
AUTHOR{{with $length := len .Metadata.Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Metadata.Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}
COMMANDS:{{range .VisibleCategories}}{{if .Name}}
   {{.Name}}:{{end}}{{range .VisibleCommands}}
     {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Metadata.Copyright}}
COPYRIGHT:
   {{.Metadata.Copyright}}{{end}}
