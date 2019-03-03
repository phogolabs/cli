# AWS S3 Provider

A package that facilitates working with http://vaultproject.io/ in context of
[CLI](https://github.com/phogolabs/cli). It increases the security of Golang
applications by populating a command line arguments from the vault.

## Installation

Make sure you have a working Go environment. Go version 1.2+ is supported.

[See the install instructions for Go](http://golang.org/doc/install.html).

To install vault, simply run:

```
$ go get github.com/phogolabs/cli/provider/aws/s3
```

## Getting Started

In order to have the provider enabled, you need to set its token either
directly or authenticating the client with Kuberenetes. For that purpose, you
will need to set the following flags in your application:

```golang
import (
	"os"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/provider/aws/s3"
)

func main() {
	app := &cli.App{
		Name:      "prana",
		HelpName:  "prana",
		Usage:     "Golang Database Manager",
		UsageText: "prana [global options]",
		Version:   "1.0-beta-04",
		Action:    run,
		Providers: []cli.Provider{
			&s3.Provider{},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "aws-access-key-id",
				Usage:  "AWS Access Key ID",
				EnvVar: "AWS_ACCESS_KEY_ID",
			},
			&cli.StringFlag{
				Name:   "aws-secret-access-key",
				Usage:  "AWS Secret Access Key",
				EnvVar: "AWS_SECRET_ACCESS_KEY",
			},
			&cli.StringFlag{
				Name:   "aws-region",
				Usage:  "AWS Region",
				EnvVar: "AWS_DEFAULT_REGION",
			},
			&cli.StringFlag{
				Name:   "aws-bucket",
				Usage:  "AWS S3 Bucket",
				EnvVar: "AWS_BUCKET",
			},
			&cli.StringFlag{
				Name:   "config",
				Usage:  "Aplication's config",
				EnvVar: "APP_CONFIG",
				FilePath: "s3://config.json",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func run(ctx *cli.Context) error {
	fmt.Println("Application started")
	return nil
}
```

As you can see in order to match the flag with a file on S# you should set
the `FilePath` field in the following format:

```
s3://<your_file_name>
```
