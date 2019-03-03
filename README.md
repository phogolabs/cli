# CLI

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][travis-img]][travis-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

A simple package for building command line applications in Go. The API is
influenced by https://github.com/urfave/cli package, but it is way more
flexible. It provides the following features:

- Data conversion that allows conversion of data to a compatible data type accepted by the declared flag
- Data providers that allows setting the flag's value from different sources such as environment variables, files, AWS S3 and SSM as well as Hashi Vault
- More extensible flag types such as URL, IP, JSON, YAML and so on.

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

	"github.com/phogolabs/cli"
)

func main() {
	app := &cli.App{
		Name:      "prana",
		HelpName:  "prana",
		Usage:     "Golang Database Manager",
		UsageText: "prana [global options]",
		Version:   "1.0-beta-04",
		Action:    run,
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

## Providers

The providers allow setting the flag's value from external sources:

- [Vault](./provider/vault/README.md) - a secure way to store and rotate credentials

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
