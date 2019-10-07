# AWS Provider

A package that contains a providers that populates flags from Amazon S3 or
Amazon SSM.

## Installation

Make sure you have a working Go environment. Go version 1.13.x is supported.

[See the install instructions for Go](http://golang.org/doc/install.html).

To install vault, simply run:

```
$ go get github.com/phogolabs/cli/provider/aws
```

## S3 Provider

This provider populates the flags from S3 bucket file.

### Getting Started

As you can see in order to match the flag with a file on S3 you should set
the `FilePath` field in the following format:

```
s3://<your_file_name>
```

```golang
import (
	"os"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/cli/provider/aws"
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
			&aws.S3Provider{},
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

	app.Run(os.Args)
}

func run(ctx *cli.Context) error {
	fmt.Println("Application started")
	return nil
}
```

## SSM Provider

This provider populates the flags from AWS SSM store.

### Getting Started

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
				Metadata: cli.Map{
				  "ssm_param": "/your/key/name",
				},
			},
		},
	}

	app.Run(os.Args)
}

func run(ctx *cli.Context) error {
	fmt.Println("Application started")
	return nil
}
```

