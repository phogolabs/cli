package cli

import (
	"bytes"
	"fmt"
	"strings"
)

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

		fmt.Fprintf(buffer, item)
	}

	return buffer.String()
}
