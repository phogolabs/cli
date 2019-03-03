package cli

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//go:generate counterfeiter -fake-name Provider -o ./fake/provider.go . Provider

// Provider is the interface that parses the flags
type Provider interface {
	Provide(*Context) error
}

// Transaction restore its state
type Transaction interface {
	Rollback(*Context) error
}

var _ Provider = &FlagProvider{}

// FlagProvider parses the CLI flags
type FlagProvider struct {
	set *flag.FlagSet
}

// Provide parses the args
func (p *FlagProvider) Provide(ctx *Context) error {
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

var _ Provider = &EnvProvider{}

// EnvProvider parses environment variables
type EnvProvider struct{}

// Provide parses the args
func (p *EnvProvider) Provide(ctx *Context) error {
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

var _ Provider = &FileProvider{}

// FileProvider parses flags from file
type FileProvider struct{}

// Provide parses the args
func (p *FileProvider) Provide(ctx *Context) error {
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
	_ Provider    = &DefaultValueProvider{}
	_ Transaction = &DefaultValueProvider{}
)

// DefaultValueProvider keeps the default values
type DefaultValueProvider struct {
	values map[string]interface{}
}

// Provide parses the args
func (p *DefaultValueProvider) Provide(ctx *Context) error {
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

// Rollback rollbacks the values
func (p *DefaultValueProvider) Rollback(ctx *Context) error {
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