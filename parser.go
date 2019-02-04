package cli

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
)

// ParserContext is the parser's context
type ParserContext struct {
	// Name of the app
	Name string
	// Args are the command line arguments
	Args []string
	// Flags are the app's defined flags
	Flags []Flag
	// Output sets the destination for usage and error messages. If output is nil, os.Stderr is used.
	Output io.Writer
}

//go:generate counterfeiter -fake-name Parser -o ./fake/parser.go . Parser

// Parser is the interface that parses the flags
type Parser interface {
	Parse(*ParserContext) error
}

var _ Parser = &FlagParser{}

// FlagParser parses the CLI flags
type FlagParser struct {
	set *flag.FlagSet
}

// Parse parses the args
func (p *FlagParser) Parse(ctx *ParserContext) error {
	p.set = flag.NewFlagSet(ctx.Name, flag.ContinueOnError)
	p.set.SetOutput(ctx.Output)

	for _, flag := range ctx.Flags {
		definition := flag.Definition()
		p.set.Var(flag, definition.Name, definition.Usage)
	}

	return p.set.Parse(ctx.Args)
}

var _ Parser = &EnvParser{}

// EnvParser parses environment variables
type EnvParser struct{}

// Parse parses the args
func (p *EnvParser) Parse(ctx *ParserContext) error {
	for _, flag := range ctx.Flags {
		definition := flag.Definition()

		env := definition.EnvVar
		if env == "" {
			continue
		}

		value := os.Getenv(env)
		if value == "" {
			continue
		}

		if err := flag.Set(value); err != nil {
			return err
		}
	}

	return nil
}

var _ Parser = &FileParser{}

// FileParser parses flags from file
type FileParser struct{}

// Parse parses the args
func (p *FileParser) Parse(ctx *ParserContext) error {
	for _, flag := range ctx.Flags {
		definition := flag.Definition()

		data, err := ioutil.ReadFile(definition.FilePath)
		if err != nil {
			continue
		}

		value := string(data)
		if err := flag.Set(value); err != nil {
			return err
		}
	}

	return nil
}
