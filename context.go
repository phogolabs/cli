package cli

import (
	"io"
	"strings"
)

// Context represents the execution context
type Context struct {
	// Args are the command line arguments
	Args []string
	// Command that owns the context
	Command *Command
	// Parent Context
	Parent *Context
	// Writer writer to write output to
	Writer io.Writer
	// ErrWriter writes error output
	ErrWriter io.Writer
}

// Bool looks up the value of a local BoolFlag, returns
// false if not found
func (ctx *Context) Bool(name string) bool {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(bool); ok {
			return value
		}
	}

	return false
}

// GlobalBool looks up the value of a global BoolFlag, returns
// false if not found
func (ctx *Context) GlobalBool(name string) bool {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(bool); ok {
			return value
		}
	}

	return false
}

// String looks up the value of a local StringFlag, returns "" if not found
func (ctx *Context) String(name string) string {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(string); ok {
			return value
		}
	}

	return ""
}

// GlobalString looks up the value of a global StringFlag, returns "" if not found
func (ctx *Context) GlobalString(name string) string {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(string); ok {
			return value
		}
	}

	return ""
}

// StringSlice looks up the value of a local StringSliceFlag, returns
// nil if not found
func (ctx *Context) StringSlice(name string) []string {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().([]string); ok {
			return value
		}
	}

	return []string{}
}

// GlobalStringSlice looks up the value of a global StringSliceFlag, returns
// nil if not found
func (ctx *Context) GlobalStringSlice(name string) []string {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().([]string); ok {
			return value
		}
	}

	return []string{}
}

func (ctx *Context) find(name string) Flag {
	for _, flag := range ctx.Command.Flags {
		names := strings.Split(flag.Definition().Name, ",")

		for _, key := range names {
			key = strings.TrimSpace(key)

			if strings.EqualFold(name, key) {
				return flag
			}
		}
	}

	return nil
}

func (ctx *Context) findAncestor(name string) Flag {
	if ctx.Parent != nil {
		ctx = ctx.Parent
	}

	for ctx != nil {
		if flag := ctx.find(name); flag != nil {
			return flag
		}

		ctx = ctx.Parent
	}

	return nil
}
