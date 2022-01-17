# CLI

[![Documentation][godoc-img]][godoc-url]
![License][license-img]
[![Build Status][action-img]][action-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

A simple package for building command line applications in Go. The API is
influenced by https://github.com/urfave/cli package, but it is way more
flexible. It provides the following features:

- More extensible flag types such as URL, IP, JSON, YAML and so on. For more information see the [docs][godoc-url]
- Data providers that allow setting the flag's value from different sources such as command line arguments, environment variables and [etc](https://github.com/hairyhenderson/go-fsimpl).
- Data conversion that allow conversion of data to a compatible data type accepted by the declared flag

## Installation

Make sure you have a working Go environment. Go version 1.16.x is supported.

[See the install instructions for Go](http://golang.org/doc/install.html).

To install CLI, simply run:

```
$ go get github.com/phogolabs/cli
```

## Getting Started

```golang
package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/phogolabs/cli"
)

var flags = []cli.Flag{
	&cli.YAMLFlag{
		Name:     "config",
		Usage:    "Application Config",
		Path:     "/etc/app/default.conf",
		EnvVar:   "APP_CONFIG",
		Value:    &Config{},
		Required: true,
	},
	&cli.StringFlag{
		Name:     "listen-addr",
		Usage:    "Application TCP Listen Address",
		EnvVar:   "APP_LISTEN_ADDR",
		Value:    ":8080",
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

As a struct that has a `Validate` function:

``` golang
type Validator struct{}

func (v *Validator) Validate(ctx *cli.Context, value interface{}) error {
	//TODO: your validation logic
	return nil
}
```

Then you can set the validator like that:

```golang
var flags = []cli.Flag{
	&cli.StringFlag{
		Name:      "name",
		EnvVar:    "APP_NAME",
		Validator: &Validator{},
	},
}
```

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
