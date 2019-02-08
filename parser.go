package cli

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//go:generate counterfeiter -fake-name Parser -o ./fake/parser.go . Parser

// Parser is the interface that parses the flags
type Parser interface {
	Parse(*Context) error
}

var _ Parser = &FlagParser{}

// FlagParser parses the CLI flags
type FlagParser struct {
	set *flag.FlagSet
}

// Parse parses the args
func (p *FlagParser) Parse(ctx *Context) error {
	p.set = flag.NewFlagSet(ctx.Command.Name, flag.ContinueOnError)
	p.set.SetOutput(ioutil.Discard)

	for _, flag := range ctx.Command.Flags {
		definition := flag.Definition()

		for _, key := range split(definition.Name) {
			key = strings.TrimSpace(key)
			p.set.Var(flag, key, definition.Usage)
		}
	}

	err := p.set.Parse(ctx.Args)
	if err != nil {
		return err
	}

	ctx.Args = p.set.Args()
	return nil
}

var _ Parser = &EnvParser{}

// EnvParser parses environment variables
type EnvParser struct{}

// Parse parses the args
func (p *EnvParser) Parse(ctx *Context) error {
	for _, flag := range ctx.Command.Flags {
		definition := flag.Definition()

		env := definition.EnvVar
		if env == "" {
			continue
		}

		for _, value := range split(os.Getenv(env)) {
			if value == "" {
				continue
			}

			if err := flag.Set(value); err != nil {
				return err
			}
		}
	}

	return nil
}

var _ Parser = &FileParser{}

// FileParser parses flags from file
type FileParser struct{}

// Parse parses the args
func (p *FileParser) Parse(ctx *Context) error {
	for _, flag := range ctx.Command.Flags {
		definition := flag.Definition()

		if definition.FilePath == "" {
			continue
		}

		paths, err := filepath.Glob(definition.FilePath)
		if err != nil {
			return err
		}

		for _, path := range paths {
			value, err := ioutil.ReadFile(path)
			if err != nil {
				continue
			}

			if err := flag.Set(string(value)); err != nil {
				return err
			}
		}
	}

	return nil
}
