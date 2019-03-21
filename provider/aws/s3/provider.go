package s3

import (
	"net/url"
	"strings"

	"github.com/phogolabs/cli"
)

var (
	_ cli.Provider = &Provider{}
)

//go:generate counterfeiter -fake-name FileSystem -o ./fake/file_system.go . FileSystem

// FileSystem represents a file system
type FileSystem interface {
	// Glob returns a list of all paths
	Glob(pattern string) ([]string, error)
	// ReadFile reads a file from the bucket
	ReadFile(path string) ([]byte, error)
}

// Provider is a parser that populates flags from AWS S3
type Provider struct {
	FileSystem FileSystem
}

// Provide parses the args
func (m *Provider) Provide(ctx *cli.Context) error {
	if err := m.init(ctx); err != nil {
		return err
	}

	if m.FileSystem == nil {
		return nil
	}

	for _, flag := range ctx.Command.Flags {
		accessor := &cli.FlagAccessor{Flag: flag}

		for _, path := range split(accessor.FilePath()) {
			paths, err := m.FileSystem.Glob(path)
			if err != nil {
				return err
			}

			for _, path := range paths {
				value, err := m.FileSystem.ReadFile(path)
				if err != nil {
					continue
				}

				if err := accessor.SetValue(string(value)); err != nil {
					return cli.FlagError(accessor, err)
				}
			}
		}
	}

	return nil
}

func (m *Provider) init(ctx *cli.Context) error {
	if m.FileSystem != nil {
		return nil
	}

	var (
		region = ctx.String("aws-region")
		bucket = ctx.String("aws-bucket")
	)

	if region == "" || bucket == "" {
		return nil
	}

	m.FileSystem = &Client{
		Config: &ClientConfig{
			Region: region,
			Bucket: bucket,
		},
	}

	return nil
}

func split(text string) []string {
	uri, err := url.Parse(text)

	if err != nil || uri.Scheme != "s3" {
		return []string{}
	}

	text = strings.TrimPrefix(text, "s3://")
	items := strings.Split(text, ",")

	for index, item := range items {
		items[index] = strings.TrimSpace(item)
	}

	return items
}
