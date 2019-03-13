# AWS SSM Provider

A package that facilitates working with http://vaultproject.io/ in context of
[CLI](https://github.com/phogolabs/cli). It increases the security of Golang
applications by populating a command line arguments from the vault.

## Installation

Make sure you have a working Go environment. Go version 1.2+ is supported.

[See the install instructions for Go](http://golang.org/doc/install.html).

To install vault, simply run:

```
$ go get github.com/phogolabs/cli/provider/aws/ssm
```

## Getting Started

As you can see in order to match the flag with a param on AWS SSM you should set
the `ssm_param` field in the meta data:


```golang
import (
	"os"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/provider/aws/ssm"
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
			&ssm.Provider{},
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
				Name:   "secret",
				Usage:  "Aplication's secret",
				Metadata: map[string]string{
				  "ssm_param": "/your/key/name",
				},
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

