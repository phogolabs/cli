package cli

import (
	"fmt"
	"strings"
)

// ErrCommandNotFound occurs when command is not found
var ErrCommandNotFound = fmt.Errorf("cli: command not found")

var (
	_ ExitCoder = &ExitError{}
	_ ExitCoder = &MultiError{}
)

// ExitCoder is the interface checked by `App` and `Command` for a custom exit
// code
type ExitCoder interface {
	error
	ExitCode() int
}

// ExitError fulfills both the builtin `error` interface and `ExitCoder`
type ExitError struct {
	exitCode int
	message  string
}

// NewExitError makes a new *ExitError
func NewExitError(message string, exitCode int) *ExitError {
	return &ExitError{
		exitCode: exitCode,
		message:  message,
	}
}

// WrapExitError wraps an error
func WrapExitError(err error, exitCode int) *ExitError {
	return NewExitError(err.Error(), exitCode)
}

// Error returns the string message, fulfilling the interface required by
// `error`
func (err *ExitError) Error() string {
	return fmt.Sprintf("%v", err.message)
}

// ExitCode returns the exit code, fulfilling the interface required by
// `ExitCoder`
func (err *ExitError) ExitCode() int {
	return err.exitCode
}

// MultiError is an error that wraps multiple errors.
type MultiError []error

// NewMultiError creates a new MultiError. Pass in one or more errors.
func NewMultiError(err ...error) *MultiError {
	errs := MultiError(err)
	return &errs
}

// Error implements the error interface.
func (err *MultiError) Error() string {
	errs := make([]string, len(*err))

	for index, item := range *err {
		errs[index] = item.Error()
	}

	return strings.Join(errs, "\n")
}

// ExitCode returns the exit code, fulfilling the interface required by
// `ExitCoder`
func (err *MultiError) ExitCode() int {
	code := 1

	for _, merr := range *err {
		if exitErr, ok := merr.(ExitCoder); ok {
			code = exitErr.ExitCode()
		}
	}

	return code
}

// AppendError appends an error
func AppendError(errx error, err error) error {
	if errx == nil {
		return err
	}

	errs, ok := errx.(*MultiError)
	if !ok {
		errs = &MultiError{}
		*errs = append(*errs, errx)
	}

	*errs = append(*errs, err)
	return errs
}
