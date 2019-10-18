package aws

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/phogolabs/cli"
)

var (
	_ cli.Provider = &S3Provider{}
)

//go:generate counterfeiter -fake-name FileSystem -o ./fake/file_system.go . FileSystem

// FileSystem represents a file system
type FileSystem interface {
	// Glob returns a list of all paths
	Glob(bucket, pattern string) ([]string, error)
	// ReadFile reads a file from the bucket
	ReadFile(bucket, path string) ([]byte, error)
}

// S3Provider is a parser that populates flags from AWS S3
type S3Provider struct {
	FileSystem FileSystem
}

// Provide parses the args
func (m *S3Provider) Provide(ctx *cli.Context) error {
	if err := m.init(ctx); err != nil {
		return err
	}

	if m.FileSystem == nil {
		return nil
	}

	bucket := ctx.String("aws-bucket")

	if bucket == "" {
		return nil
	}

	for _, flag := range ctx.Command.Flags {
		accessor := cli.NewFlagAccessor(flag)

		for _, path := range m.split(accessor.FilePath()) {
			paths, err := m.FileSystem.Glob(bucket, path)
			if err != nil {
				return err
			}

			for _, path := range paths {
				value, err := m.FileSystem.ReadFile(bucket, path)
				if err != nil {
					continue
				}

				if err := accessor.Set(string(value)); err != nil {
					return cli.FlagError("s3", accessor.Name(), err)
				}
			}
		}
	}

	return nil
}

func (m *S3Provider) init(ctx *cli.Context) error {
	if m.FileSystem != nil {
		return nil
	}

	var (
		region = ctx.String("aws-region")
		role   = ctx.String("aws-role-arn")
	)

	if region == "" {
		return nil
	}

	m.FileSystem = &Client{
		Config: &ClientConfig{
			Region:  region,
			RoleARN: role,
		},
	}

	return nil
}

func (m *S3Provider) split(text string) []string {
	var (
		result []string
		items  = strings.Split(text, ",")
	)

	for _, item := range items {
		item = strings.TrimSpace(item)

		uri, err := url.Parse(item)
		if err != nil {
			continue
		}

		if uri.Scheme == "s3" {
			item = strings.TrimPrefix(item, "s3://")
			result = append(result, item)
		}
	}

	return result
}

//go:generate counterfeiter -fake-name Client -o ./fake/client.go . Getter

// Getter represent the client that fetches ssm
type Getter interface {
	Get(name string) (string, error)
}

// SSMProvider is a parser that populates flags from AWS S3
type SSMProvider struct {
	Client Getter
}

// Provide parses the args
func (m *SSMProvider) Provide(ctx *cli.Context) error {
	if err := m.init(ctx); err != nil {
		return err
	}

	if m.Client == nil {
		return nil
	}

	for _, flag := range ctx.Command.Flags {
		accessor := cli.NewFlagAccessor(flag)

		meta := accessor.MetaKey("ssm_param")
		if meta == nil {
			continue
		}

		path := fmt.Sprintf("%v", meta)

		value, err := m.Client.Get(path)
		if err != nil {
			return err
		}

		if err := accessor.Set(value); err != nil {
			return cli.FlagError("ssm", accessor.Name(), err)
		}
	}

	return nil
}

func (m *SSMProvider) init(ctx *cli.Context) error {
	if m.Client != nil {
		return nil
	}

	var (
		region = ctx.String("aws-region")
		role   = ctx.String("aws-role-arn")
	)

	if region == "" {
		return nil
	}

	m.Client = &Client{
		Config: &ClientConfig{
			Region:  region,
			RoleARN: role,
		},
	}

	return nil
}
