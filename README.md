# CLI

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][action-img]][action-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

A simple package for building command line applications in Go. The API is
influenced by https://github.com/urfave/cli package, but it is way more
flexible. It provides the following features:

- Data conversion that allow conversion of data to a compatible data type accepted by the declared flag
- Data providers that allow setting the flag's value from different sources such as environment variables, files, AWS S3 and SSM
- More extensible flag types such as URL, IP, JSON, YAML and so on. For more information see the [docs][godoc-url]

## Installation

Make sure you have a working Go environment. Go version 1.13.x is supported.

[See the install instructions for Go](http://golang.org/doc/install.html).

To install CLI, simply run:

```
$ go get github.com/phogolabs/cli
```

## Getting Started

```golang
import (
	"os"
	"syscall"

	"github.com/phogolabs/cli"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:     "config",
		Usage:    "Application Config",
		EnvVar:   "APP_CONFIG",
		FilePath: "/etc/app/app.config",
		Required: true,
	},
}


func main() {
	app := &cli.App{
		Name:      "prana",
		HelpName:  "prana",
		Usage:     "Golang Database Manager",
		UsageText: "prana [global options]",
		Version:   "1.0-beta-04",
		Flags:     flags,
		Action:    run,
		Signals:   []os.Signal{syscall.SIGTERM},
		OnSignal:  signal,
	}

	app.Run(os.Args)
}

// run executes the application
func run(ctx *cli.Context) error {
	fmt.Println("Application started")
	return nil
}

// signal handles OS signal
func signal(ctx *cli.Context, signal os.Signal) error {
	fmt.Println("Application signal", signal)
	return nil
}
```

## Validation

You can set the `Required` field to `true` if you want to make some flags
mandatory. If you need some customized validation, you can create a custom
validator in the following way:

As a function:

```golang
validate := cli.ValidatorFunc(func(ctx *cli.Context, value interface{}) error {
        //TODO: your validation logic
	return nil
})
```

As a struct that has a `Validate` function:

``` golang
type Validator struct {}

func (v *Validator) Validate(ctx *cli.Context, value interface{}) error {
        //TODO: your validation logic
	return nil
}
```

Then you can set the validator like that:

```golang
var flags = []cli.Flag{
	&cli.StringSliceFlag{
		Name:     "endpoint",
		Usage:    "Endpoint URL",
		EnvVar:   "APP_ENDPOINT",
		Required: true,
		Validator: validate,
	},
	&cli.StringFlag{
		Name:      "name",
		EnvVar:    "APP_NAME",
		Validator: &Validator{},
	},
}
```

## Providers

The providers allow setting the flag's value from external sources:

- [AWS S3 Provider](./provider/aws/README.md#s3-provider) - reads a flag's value from AWS S3 bucket object
- [AWS SSM Provider](./provider/aws/README.md#ssm-provider) - reads a flag's value from AWS SSM Parameter store

## Contributing

We are open for any contributions. Just fork the
[project](https://github.com/phogolabs/cli).

[report-img]: https://goreportcard.com/badge/github.com/phogolabs/cli
[report-url]: https://goreportcard.com/report/github.com/phogolabs/cli
[codecov-url]: https://codecov.io/gh/phogolabs/cli
[codecov-img]: https://codecov.io/gh/phogolabs/cli/branch/master/graph/badge.svg
[action-img]: https://github.com/phogolabs/cli/workflows/main/badge.svg
[action-url]: https://github.com/phogolabs/cli/actions
[godoc-url]: https://godoc.org/github.com/phogolabs/cli
[godoc-img]: https://godoc.org/github.com/phogolabs/cli?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
