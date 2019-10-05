NAME:
   {{.HelpName}} - {{.Usage}}
USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}}{{if .Metadata.VisibleFlags}} [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Category}}
CATEGORY:
   {{.Category}}{{end}}{{if .Description}}
DESCRIPTION:
   {{.Description}}{{end}}{{if .Metadata.VisibleFlags}}
OPTIONS:
   {{range .Metadata.VisibleFlags}}{{.}}
   {{end}}{{end}}
