# CLI

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][travis-img]][travis-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

A simple package for building command line applications in Go. The API is
influenced by https://github.com/urfave/cli package, but it is way more
flexible. It provides the following features:

- Data conversion that allow conversion of data to a compatible data type accepted by the declared flag
- Data providers that allow setting the flag's value from different sources such as environment variables, files, AWS S3 and SSM as well as Hashi Vault
- More extensible flag types such as URL, IP, JSON, YAML and so on. For more information see the [docs][godoc-url]

## Installation

Make sure you have a working Go environment. Go version 1.2+ is supported.

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
		Name:   "aws-region",
		Usage:  "AWS Region",
		EnvVar: "AWS_REGION, AWS_DEFAULT_REGION",
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
		Signals:   []os.Signal{system.SIGTERM},
		OnSignal:  signal,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

// run executes the application
func run(ctx *cli.Context) error {
	fmt.Println("Application started")
	return nil
}

// signal handles OS signal
func signal(ctx *cli.Context) error {
	fmt.Println("Application signal", ctx.Signal)
	return nil
}
```

## Providers

The providers allow setting the flag's value from external sources:

- [Vault](https://github.com/phogolabs/vault/blob/master/docs/provider.md) - reads a flag's value from Hashi Corp Vault Secret
- [AWS S3](./provider/aws/s3/README.md) - reads a flag's value from AWS S3 bucket object
- [AWS SSM](./provider/aws/ssm/README.md) - reads a flag's value from AWS SSM Parameter store

## Converters

Let's assume that we have the following JSON in your KV config:

```json
{
  "username": "root",
  "password": "swordfish"
}
```

If you want to populate a flag's value with the password field you should use
[JSON Path](https://goessner.net/articles/JsonPath/) by setting the flag's
converter to `cli.JSONPath`:

```golang
flag := &cli.StringFlag{
	Name:   "password",
	Usage:  "Aplication's password",
	FilePath: "app.config",
	Converter: cli.JSONPath("$.password"),
}
```

## Contributing

We are welcome to any contributions. Just fork the
[project](https://github.com/phogolabs/cli).

[travis-img]: https://travis-ci.org/phogolabs/cli.svg?branch=master
[travis-url]: https://travis-ci.org/phogolabs/cli
[report-img]: https://goreportcard.com/badge/github.com/phogolabs/cli
[report-url]: https://goreportcard.com/report/github.com/phogolabs/cli
[codecov-url]: https://codecov.io/gh/phogolabs/cli
[codecov-img]: https://codecov.io/gh/phogolabs/cli/branch/master/graph/badge.svg
[godoc-url]: https://godoc.org/github.com/phogolabs/cli
[godoc-img]: https://godoc.org/github.com/phogolabs/cli?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[software-license-url]: LICENSE
