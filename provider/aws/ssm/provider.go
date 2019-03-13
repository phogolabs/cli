package ssm

import (
	"github.com/phogolabs/cli"
)

//go:generate counterfeiter -fake-name Client -o ./fake/client.go . Getter

// Getter represent the client that fetches ssm
type Getter interface {
	Get(name string) (string, error)
}

// Provider is a parser that populates flags from AWS S3
type Provider struct {
	Client Getter
}

// Provide parses the args
func (m *Provider) Provide(ctx *cli.Context) error {
	if err := m.init(ctx); err != nil {
		return err
	}

	if m.Client == nil {
		return nil
	}

	for _, flag := range ctx.Command.Flags {
		accessor := &cli.FlagAccessor{Flag: flag}

		path := accessor.MetaKey("ssm_param")
		if path == "" {
			continue
		}

		value, err := m.Client.Get(path)
		if err != nil {
			return err
		}

		if err := accessor.SetValue(value); err != nil {
			return err
		}
	}

	return nil
}

func (m *Provider) init(ctx *cli.Context) error {
	if m.Client != nil {
		return nil
	}

	var (
		region = ctx.String("aws-region")
	)

	if region == "" {
		return nil
	}

	m.Client = &Client{
		Config: &ClientConfig{
			Region: region,
		},
	}

	return nil
}
