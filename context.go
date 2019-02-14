package cli

import (
	"io"
	"net"
	"net/url"
	"strings"
	"time"
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

// URL looks up the value of a local URLFlag, returns nil if not found
func (ctx *Context) URL(name string) *url.URL {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(*url.URL); ok {
			return value
		}
	}

	return nil
}

// GlobalURL looks up the value of a global URLFlag, returns nil if not found
func (ctx *Context) GlobalURL(name string) *url.URL {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(*url.URL); ok {
			return value
		}
	}

	return nil
}

// Time looks up the value of a local TimeFlag, returns 0 if not found
func (ctx *Context) Time(name string) time.Time {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(time.Time); ok {
			return value
		}
	}

	return time.Time{}
}

// GlobalTime looks up the value of a global TimeFlag, returns 0 if not found
func (ctx *Context) GlobalTime(name string) time.Time {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(time.Time); ok {
			return value
		}
	}

	return time.Time{}
}

// Duration looks up the value of a local DurationFlag, returns 0 if not found
func (ctx *Context) Duration(name string) time.Duration {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(time.Duration); ok {
			return value
		}
	}

	return time.Duration(0)
}

// GlobalDuration looks up the value of a global DurationFlag, returns 0 if not found
func (ctx *Context) GlobalDuration(name string) time.Duration {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(time.Duration); ok {
			return value
		}
	}

	return time.Duration(0)
}

// Float32 looks up the value of a local Float32Flag, returns 0 if not found
func (ctx *Context) Float32(name string) float32 {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(float32); ok {
			return value
		}
	}

	return 0
}

// GlobalFloat32 looks up the value of a global Float64Flag, returns 0 if not found
func (ctx *Context) GlobalFloat32(name string) float32 {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(float32); ok {
			return value
		}
	}

	return 0
}

// Float64 looks up the value of a local Float64Flag, returns 0 if not found
func (ctx *Context) Float64(name string) float64 {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(float64); ok {
			return value
		}
	}

	return 0
}

// GlobalFloat64 looks up the value of a global Float64Flag, returns 0 if not found
func (ctx *Context) GlobalFloat64(name string) float64 {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(float64); ok {
			return value
		}
	}

	return 0
}

// Int looks up the value of a local IntFlag, returns 0 if not found
func (ctx *Context) Int(name string) int {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(int); ok {
			return value
		}
	}

	return 0
}

// GlobalInt looks up the value of a global IntFlag, returns 0 if not found
func (ctx *Context) GlobalInt(name string) int {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(int); ok {
			return value
		}
	}

	return 0
}

// Int64 looks up the value of a local Int64Flag, returns 0 if not found
func (ctx *Context) Int64(name string) int64 {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(int64); ok {
			return value
		}
	}

	return 0
}

// GlobalInt64 looks up the value of a global Int64Flag, returns 0 if not found
func (ctx *Context) GlobalInt64(name string) int64 {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(int64); ok {
			return value
		}
	}

	return 0
}

// UInt looks up the value of a local UIntFlag, returns 0 if not found
func (ctx *Context) UInt(name string) uint {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(uint); ok {
			return value
		}
	}

	return 0
}

// GlobalUInt looks up the value of a global UIntFlag, returns 0 if not found
func (ctx *Context) GlobalUInt(name string) uint {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(uint); ok {
			return value
		}
	}

	return 0
}

// UInt64 looks up the value of a local UInt64Flag, returns 0 if not found
func (ctx *Context) UInt64(name string) uint64 {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(uint64); ok {
			return value
		}
	}

	return 0
}

// GlobalUInt64 looks up the value of a global UInt64Flag, returns 0 if not found
func (ctx *Context) GlobalUInt64(name string) uint64 {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(uint64); ok {
			return value
		}
	}

	return 0
}

// IP looks up the value of a local IPFlag, returns nil if not found
func (ctx *Context) IP(name string) net.IP {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(net.IP); ok {
			return value
		}
	}

	return nil
}

// GlobalIP looks up the value of a global IPFlag, returns nil if not found
func (ctx *Context) GlobalIP(name string) net.IP {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(net.IP); ok {
			return value
		}
	}

	return nil
}

// HardwareAddr looks up the value of a local HardwareddrFlag, returns nil if not found
func (ctx *Context) HardwareAddr(name string) net.HardwareAddr {
	if flag := ctx.find(name); flag != nil {
		if value, ok := flag.Get().(net.HardwareAddr); ok {
			return value
		}
	}

	return nil
}

// GlobalHardwareAddr looks up the value of a global HardwareAddrFlag, returns nil if not found
func (ctx *Context) GlobalHardwareAddr(name string) net.HardwareAddr {
	if flag := ctx.findAncestor(name); flag != nil {
		if value, ok := flag.Get().(net.HardwareAddr); ok {
			return value
		}
	}

	return nil
}

// Get looks up the value of a local flag, returns nil if not found
func (ctx *Context) Get(name string) interface{} {
	if flag := ctx.find(name); flag != nil {
		return flag.Get()
	}

	return nil
}

// GlobalGet looks up the value of a global flag, returns nil if not found
func (ctx *Context) GlobalGet(name string) interface{} {
	if flag := ctx.findAncestor(name); flag != nil {
		return flag.Get()
	}

	return nil
}

func (ctx *Context) find(name string) Flag {
	for _, flag := range ctx.Command.Flags {
		accessor := &FlagAccessor{Flag: flag}
		names := strings.Split(accessor.Name(), ",")

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
