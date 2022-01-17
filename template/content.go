package template

import (
	"embed"
	"io/ioutil"
	"strings"
	"text/template"
)

//go:embed *.tpl
var content embed.FS

// Open opens the template
func Open(name string) (*template.Template, error) {
	file, err := content.Open(name)
	if err != nil {
		return nil, err
	}

	kv := template.FuncMap{
		"join": strings.Join,
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return template.New(name).Funcs(kv).Parse(string(data))
}
