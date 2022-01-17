package cli

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"
)

// BeforeFunc is an action to execute before any subcommands are run, but after
// the context is ready if a non-nil error is returned, no subcommands are run
type BeforeFunc func(*Context) error

// AfterFunc is an action to execute after any subcommands are run, but after the
// subcommand has finished it is run even if Action() panics
type AfterFunc func(*Context) error

// ActionFunc is the action to execute when no subcommands are specified
type ActionFunc func(*Context) error

// SignalFunc is an action to execute after a system signal
type SignalFunc func(*Context, os.Signal) error

// OnUsageErrorFunc is executed if an usage error occurs. This is useful for displaying
// customized usage error messages.  This function is able to replace the
// original error messages.  If this function is not set, the "Incorrect usage"
// is displayed and the execution is interrupted.
type UsageErrorFunc func(context *Context, err error) error

// OnCommandNotFoundFunc is executed if the proper command cannot be found
type CommandNotFoundFunc func(*Context, string)

// OnExitErrorHandlerFunc is executed if provided in order to handle ExitError
// values returned by Actions and Before/After functions.
type ExitErrorHandlerFunc func(err error) error

// ExitFunc is an exit function
type ExitFunc func(code int)

// OneOf returns a validator that expectes the flag value to matches one of the
// provided values
func OneOf(items ...interface{}) Validator {
	fn := func(ctx *Context, value interface{}) error {
		for _, item := range items {
			if item == value {
				return nil
			}
		}

		return fmt.Errorf("unsupported value: %v", value)
	}

	return ValidatorFunc(fn)
}

// EnvOf formats a list of environment variables
func EnvOf(items ...string) string {
	buffer := &bytes.Buffer{}

	for index, item := range items {
		if index > 0 {
			fmt.Fprintf(buffer, ", ")
		}

		item = strings.TrimSpace(item)
		item = strings.ToUpper(item)

		fmt.Fprint(buffer, item)
	}

	return buffer.String()
}

// taken from https://github.com/urfave/cli/blob/master/sort_test.go
func less(i, j string) bool {
	iRunes := []rune(i)
	jRunes := []rune(j)

	lenShared := len(iRunes)
	if lenShared > len(jRunes) {
		lenShared = len(jRunes)
	}

	for index := 0; index < lenShared; index++ {
		ir := iRunes[index]
		jr := jRunes[index]

		if lir, ljr := unicode.ToLower(ir), unicode.ToLower(jr); lir != ljr {
			return lir < ljr
		}

		if ir != jr {
			return ir < jr
		}
	}

	return i < j
}
