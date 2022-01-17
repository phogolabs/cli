package cli

import (
	"flag"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/hairyhenderson/go-fsimpl/autofs"
	"github.com/phogolabs/log"
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

var _ Provider = &PathProvider{}

// PathProvider parses flags from file
type PathProvider struct {
	IsPathFlag bool
}

// Provide parses the args
func (p *PathProvider) Provide(ctx *Context) error {
	for _, flag := range ctx.Command.Flags {
		accessor := NewFlagAccessor(flag)

		if p.IsPathFlag != accessor.IsPathFlag() {
			continue
		}

		for _, path := range split(accessor.Path()) {
			if path == "" {
				continue
			}

			root, err := p.root(path)
			if err != nil {
				return err
			}

			fs, err := autofs.Lookup(root.String())
			if err != nil {
				return err
			}

			name, err := p.name(path)
			if err != nil {
				return err
			}

			source, err := fs.Open(name)
			if err != nil {
				return err
			}
			// close the source
			defer source.Close()

			if _, err := accessor.ReadFrom(source); err != nil {
				return err
			}
		}
	}

	return nil
}

func (p *PathProvider) root(path string) (*url.URL, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	if uri.Scheme == "" {
		uri.Scheme = "file"
	}

	if dir := filepath.Dir(uri.Path); dir == "." {
		uri.Path, _ = os.Getwd()
	} else {
		uri.Path = dir
	}

	if dir := filepath.Dir(uri.RawPath); dir == "." {
		uri.RawPath, _ = os.Getwd()
	} else {
		uri.RawPath = dir
	}

	if uri.Host == "." {
		uri.Path, _ = os.Getwd()
		uri.RawPath, _ = os.Getwd()
	}

	return uri, nil
}

func (p *PathProvider) name(path string) (string, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	if name := filepath.Base(uri.Path); name != "." {
		path = name
	}

	return path, nil
}

// BackOffStrategy represents the backoff strategy
type BackOffStrategy backoff.BackOff

// BackOffProvider backoff the provider
type BackOffProvider struct {
	Provider Provider
	Strategy BackOffStrategy
}

// Provide parses the args
func (m *BackOffProvider) Provide(ctx *Context) error {
	tryProvide := func() error {
		log.Info("providing the application arguments")
		return m.Provider.Provide(ctx)
	}

	notify := func(err error, t time.Duration) {
		log.WithError(err).Warnf("providing the application arguments not successful. retry in %v", t)
	}

	if m.Strategy == nil {
		// create the default strategy
		strategy := backoff.NewExponentialBackOff()
		strategy.MaxElapsedTime = 30 * time.Second
		strategy.InitialInterval = 2 * time.Second
		// set the default strategy
		m.Strategy = strategy
	}

	if err := backoff.RetryNotify(tryProvide, m.Strategy, notify); err != nil {
		log.WithError(err).Fatal("providing the application argument failed")
		return err
	}

	return nil
}
