package cli

import (
	"fmt"
	"strings"
)

const (
	// ExitCodeErrorApp is the exit code on application error
	ExitCodeErrorApp = 1001
	// ExitCodeErrorProvider is the exit code on provider error
	ExitCodeErrorProvider = 1002
	// ExitCodeNotFoundFlag is the exit code when a flag is not found
	ExitCodeNotFoundFlag = 1003
	// ExitCodeNotFoundCommand is the exit code when a command is not found
	ExitCodeNotFoundCommand = 1004
)

// ExitCoder is the interface checked by `App` and `Command` for a custom exit
// code
type ExitCoder interface {
	Error() string
	ExitCode() int
}

var _ ExitCoder = &ExitError{}

// ExitError fulfills both the builtin `error` interface and `ExitCoder`
type ExitError struct {
	code int
	err  error
}

// NewExitError makes a new ExitError
func NewExitError(text string, code int) *ExitError {
	return &ExitError{
		code: code,
		err:  fmt.Errorf(text),
	}
}

// WrapError wraps an error as ExitError
func WrapError(err error, code int) *ExitError {
	return &ExitError{
		code: code,
		err:  err,
	}
}

// NotFoundFlagError makes a new ExitError for missing flags
func NotFoundFlagError(name string) *ExitError {
	return &ExitError{
		code: ExitCodeNotFoundFlag,
		err:  fmt.Errorf("flag '%s' not found", name),
	}
}

// NotFoundCommandError makes a new ExitError for missing command
func NotFoundCommandError(name string) *ExitError {
	return &ExitError{
		code: ExitCodeNotFoundCommand,
		err:  fmt.Errorf("command '%s' not found", name),
	}
}

// ProviderCommandError makes a new ExitError for missing command
func ProviderFlagError(name string, flag *FlagAccessor, err error) *ExitError {
	return &ExitError{
		code: ExitCodeErrorProvider,
		err:  fmt.Errorf("provider '%s' failed to set a flag '%v': %w", name, flag.Name(), err),
	}
}

// Error returns the string message, fulfilling the interface required by
// `error`
func (x *ExitError) Error() string {
	return x.err.Error()
}

// ExitCode returns the exit code, fulfilling the interface required by
// `ExitCoder`
func (x *ExitError) ExitCode() int {
	return x.code
}

// Unwrap returns the underlying error
func (x *ExitError) Unwrap() error {
	return x.err
}

var _ ExitCoder = &ExitErrorCollector{}

// ExitErrorCollector is an error that wraps multiple errors.
type ExitErrorCollector []error

// Error implements the error interface.
func (errs ExitErrorCollector) Error() string {
	messages := make([]string, len(errs))

	for index, err := range errs {
		messages[index] = err.Error()
	}

	return strings.Join(messages, "\n")
}

// ExitCode returns the exit code, fulfilling the interface required by
// ExitCoder
func (errs ExitErrorCollector) ExitCode() int {
	for _, err := range errs {
		if errx, ok := err.(ExitCoder); ok {
			return errx.ExitCode()
		}
	}

	return ExitCodeErrorApp
}

// Unwrap unwraps the error
func (errs ExitErrorCollector) Unwrap() error {
	count := len(errs)

	switch {
	case count == 1:
		return errs[0]
	default:
		return nil
	}
}
