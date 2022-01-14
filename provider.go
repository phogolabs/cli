package cli

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Map of key value pairs
type Map map[string]interface{}

//go:generate counterfeiter -fake-name Provider -o ./fake/provider.go . Provider

// Provider is the interface that parses the flags
type Provider interface {
	Provide(*Context) error
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
		accessor := NewFlagAccessor(flag)

		for _, key := range split(accessor.Name()) {
			p.set.Var(accessor, key, accessor.Usage())
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
	for _, flag := range ctx.Command.Flags {
		accessor := NewFlagAccessor(flag)

		for _, env := range split(accessor.EnvVar()) {
			value := getEnv(env)

			for _, value := range split(value) {
				if err := accessor.Set(value); err != nil {
					return FlagError("env", accessor.Name(), err)
				}
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
		accessor := NewFlagAccessor(flag)

		for _, fpath := range split(accessor.FilePath()) {
			paths, err := filepath.Glob(fpath)

			if err != nil {
				return err
			}

			for _, root := range paths {
				info, err := os.Stat(root)
				if err != nil {
					continue
				}

				var values []string

				if info.IsDir() {
					values, err = readDir(root)
					if err != nil {
						continue
					}
				} else {
					values, err = readFile(fpath)
					if err != nil {
						continue
					}
				}

				for _, value := range values {
					if err := accessor.Set(value); err != nil {
						return FlagError("file", accessor.Name(), err)
					}
				}
			}
		}
	}

	return nil
}
