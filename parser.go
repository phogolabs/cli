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

// Restorer restore its state
type Restorer interface {
	Restore(*Context) error
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
		accessor := &FlagAccessor{Flag: flag}

		for _, key := range split(accessor.Name()) {
			key = strings.TrimSpace(key)
			p.set.Var(flag, key, accessor.Usage())
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
	var env string

	for _, flag := range ctx.Command.Flags {
		accessor := &FlagAccessor{Flag: flag}

		if env = accessor.EnvVar(); env == "" {
			continue
		}

		for _, value := range split(os.Getenv(env)) {
			if value == "" {
				continue
			}

			if err := accessor.SetValue(value); err != nil {
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
		accessor := &FlagAccessor{Flag: flag}

		for _, path := range split(accessor.FilePath()) {
			paths, err := filepath.Glob(path)
			if err != nil {
				return err
			}

			for _, path := range paths {
				value, err := ioutil.ReadFile(path)
				if err != nil {
					continue
				}

				if err := accessor.SetValue(string(value)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

var (
	_ Parser   = &DefaultValueParser{}
	_ Restorer = &DefaultValueParser{}
)

// DefaultValueParser keeps the default values
type DefaultValueParser struct {
	values map[string]interface{}
}

// Parse parses the args
func (p *DefaultValueParser) Parse(ctx *Context) error {
	if p.values == nil {
		p.values = make(map[string]interface{})
	}

	for _, flag := range ctx.Command.Flags {
		accessor := &FlagAccessor{Flag: flag}

		if value := accessor.Value(); value != nil {
			p.values[accessor.Name()] = value
		}
	}

	return nil
}

// Restore rollbacks the values
func (p *DefaultValueParser) Restore(ctx *Context) error {
	for _, flag := range ctx.Command.Flags {
		accessor := &FlagAccessor{Flag: flag}

		if value, ok := p.values[accessor.Name()]; ok {
			if err := accessor.SetValue(value); err != nil {
				return err
			}
		}
	}

	return nil
}
