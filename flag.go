package cli

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"gopkg.in/yaml.v2"
)

//go:generate counterfeiter -fake-name Flag -o ./fake/flag.go . Flag

// Flag is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
//
// If a Value has an IsBoolFlag() bool method returning true,
// the command-line parser makes -name equivalent to -name=true
// rather than using the next command-line argument.
//
// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
//
// Getter is an interface that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
type Flag interface {
	String() string
	Set(string) error
	Get() interface{}
}

var _ Flag = &StringFlag{}

// StringFlag is a flag with type string
type StringFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     string
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *StringFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *StringFlag) Set(value string) error {
	f.Value = value
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *StringFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *StringFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == "" {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &StringSliceFlag{}

// StringSliceFlag is a flag with type *StringSlice
type StringSliceFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     []string
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *StringSliceFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *StringSliceFlag) Set(value string) error {
	f.Value = append(f.Value, value)
	return nil
}

// Reset resets the valle
func (f *StringSliceFlag) Reset() error {
	f.Value = []string{}
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *StringSliceFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *StringSliceFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == nil || len(f.Value) == 0 {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &BoolFlag{}

// BoolFlag is a flag with type bool
type BoolFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     bool
	Hidden    bool
	Validator Validator
}

// IsBoolFlag returns true if the flag is bool
func (f *BoolFlag) IsBoolFlag() bool {
	return true
}

// String returns the value as string
func (f *BoolFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *BoolFlag) Set(value string) error {
	flag, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}

	f.Value = flag
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *BoolFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *BoolFlag) Validate(ctx *Context) error {
	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &URLFlag{}

// URLFlag is a flag with type url.URL
type URLFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     *url.URL
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *URLFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *URLFlag) Set(value string) error {
	uri, err := url.Parse(value)
	if err != nil {
		return err
	}

	f.Value = uri
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *URLFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *URLFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == nil {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &JSONFlag{}

// JSONFlag is a flag with type json document
type JSONFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     interface{}
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *JSONFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *JSONFlag) Set(value string) error {
	f.Path = value
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *JSONFlag) Get() interface{} {
	return f.Path
}

// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except EOF encountered during the read is also returned.
func (f *JSONFlag) ReadFrom(r io.Reader) (int64, error) {
	if f.Value == nil {
		f.Value = &map[string]interface{}{}
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	n := len(data)
	// decode from base64
	if content, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
		data = content
	}

	if err := json.Unmarshal(data, f.Value); err != nil {
		return 0, err
	}

	return int64(n), nil
}

// Validate validates the flag
func (f *JSONFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Path == "" || f.Value == nil {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

// IsPathFlag returns true if the flag is path
func (f *JSONFlag) IsPathFlag() bool {
	return true
}

var _ Flag = &YAMLFlag{}

// YAMLFlag is a flag with type yaml document
type YAMLFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     interface{}
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *YAMLFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *YAMLFlag) Set(value string) error {
	f.Path = value
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *YAMLFlag) Get() interface{} {
	return f.Path
}

// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except EOF encountered during the read is also returned.
func (f *YAMLFlag) ReadFrom(r io.Reader) (int64, error) {
	if f.Value == nil {
		f.Value = &map[string]interface{}{}
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	n := len(data)
	// decode from base64
	if content, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
		data = content
	}

	if err := yaml.Unmarshal(data, f.Value); err != nil {
		return 0, err
	}

	return int64(n), nil
}

// Validate validates the flag
func (f *YAMLFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Path == "" || f.Value == nil {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

// IsPathFlag returns true if the flag is path
func (f *YAMLFlag) IsPathFlag() bool {
	return true
}

var _ Flag = &XMLFlag{}

// XMLFlag is a flag with type XMLDocument
type XMLFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     interface{}
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *XMLFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *XMLFlag) Set(value string) error {
	f.Path = value
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *XMLFlag) Get() interface{} {
	return f.Path
}

// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except EOF encountered during the read is also returned.
func (f *XMLFlag) ReadFrom(r io.Reader) (int64, error) {
	if f.Value == nil {
		f.Value = &map[string]interface{}{}
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	n := len(data)
	// decode from base64
	if content, err := base64.StdEncoding.DecodeString(string(data)); err == nil {
		data = content
	}

	if err := xml.Unmarshal(data, f.Value); err != nil {
		return 0, err
	}

	return int64(n), err
}

// Validate validates the flag
func (f *XMLFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Path == "" || f.Value == nil {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

// IsPathFlag returns true if the flag is path
func (f *XMLFlag) IsPathFlag() bool {
	return true
}

var _ Flag = &TimeFlag{}

// TimeFlag is a flag with type time.Time
type TimeFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Format    string
	Value     time.Time
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *TimeFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *TimeFlag) Set(value string) (err error) {
	if f.Format == "" {
		f.Format = time.UnixDate
	}

	f.Value, err = time.Parse(f.Format, value)
	return
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *TimeFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *TimeFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value.IsZero() {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &DurationFlag{}

// DurationFlag is a flag with type time.Duration
type DurationFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     time.Duration
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *DurationFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *DurationFlag) Set(value string) (err error) {
	f.Value, err = time.ParseDuration(value)
	return
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *DurationFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *DurationFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == 0 {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &IntFlag{}

// IntFlag is a flag with type int
type IntFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     int
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *IntFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *IntFlag) Set(value string) error {
	parsed, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		return err
	}

	f.Value = int(parsed)
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *IntFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *IntFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == 0 {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &Int64Flag{}

// Int64Flag is a flag with type int64
type Int64Flag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     int64
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *Int64Flag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *Int64Flag) Set(value string) (err error) {
	f.Value, err = strconv.ParseInt(value, 0, 64)
	return
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *Int64Flag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *Int64Flag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == 0 {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &UIntFlag{}

// UIntFlag is a flag with type uint64
type UIntFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     uint
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *UIntFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *UIntFlag) Set(value string) error {
	parsed, err := strconv.ParseUint(value, 0, 64)
	if err != nil {
		return err
	}

	f.Value = uint(parsed)
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *UIntFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *UIntFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == 0 {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &UInt64Flag{}

// UInt64Flag is a flag with type uint
type UInt64Flag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     uint64
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *UInt64Flag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *UInt64Flag) Set(value string) (err error) {
	f.Value, err = strconv.ParseUint(value, 0, 64)
	return
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *UInt64Flag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *UInt64Flag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == 0 {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &Float32Flag{}

// Float32Flag is a flag with type float32
type Float32Flag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     float32
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *Float32Flag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *Float32Flag) Set(value string) error {
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	f.Value = float32(parsed)
	return nil
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *Float32Flag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *Float32Flag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == 0 {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &Float64Flag{}

// Float64Flag is a flag with type float64
type Float64Flag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     float64
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *Float64Flag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *Float64Flag) Set(value string) (err error) {
	f.Value, err = strconv.ParseFloat(value, 64)
	return
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *Float64Flag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *Float64Flag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == 0 {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &IPFlag{}

// IPFlag is a flag with type net.IP
type IPFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     net.IP
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *IPFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *IPFlag) Set(value string) (err error) {
	f.Value = net.ParseIP(value)

	if f.Value == nil && value != "" {
		return &net.ParseError{
			Type: "IP Address",
			Text: value,
		}
	}

	return
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *IPFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *IPFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == nil {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &HardwareAddrFlag{}

// HardwareAddrFlag is a flag with type net.HardwareAddr
type HardwareAddrFlag struct {
	Name      string
	Path      string
	Usage     string
	EnvVar    string
	Value     net.HardwareAddr
	Hidden    bool
	Required  bool
	Validator Validator
}

// String returns the value as string
func (f *HardwareAddrFlag) String() string {
	return FlagFormat(f)
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *HardwareAddrFlag) Set(value string) (err error) {
	f.Value, err = net.ParseMAC(value)
	return
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *HardwareAddrFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *HardwareAddrFlag) Validate(ctx *Context) error {
	if f.Required {
		if f.Value == nil {
			return NotFoundFlagError(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(ctx, f.Value)
	}

	return nil
}

var _ Flag = &FlagAccessor{}

// FlagAccessor access the flag's field
type FlagAccessor struct {
	Flag  Flag
	Text  string
	IsSet bool
}

// NewFlagAccessor returns new flag accessor
func NewFlagAccessor(flag Flag) *FlagAccessor {
	return &FlagAccessor{
		Flag: flag,
		Text: flag.String(),
	}
}

// String returns the flag as string
func (f *FlagAccessor) String() string {
	return f.Text
}

// Set is called once, in command line order, for each flag present.
// The flag package may call the String method with a zero-valued receiver,
// such as a nil pointer.
func (f *FlagAccessor) Set(value string) error {
	if value == "" {
		return nil
	}

	if !f.IsSet {
		// reset the value
		if err := f.Reset(); err != nil {
			return err
		}

		f.IsSet = true
	}

	return f.Flag.Set(value)
}

// Get of the flag's value
func (f *FlagAccessor) Get() interface{} {
	return f.Flag.Get()
}

// ReadFrom reads data from r until EOF or error.
// The return value n is the number of bytes read.
// Any error except EOF encountered during the read is also returned.
func (f *FlagAccessor) ReadFrom(r io.Reader) (int64, error) {
	if reader, ok := f.Flag.(io.ReaderFrom); ok {
		return reader.ReadFrom(r)
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	n := len(data)
	// use the set method
	if err := f.Set(string(data)); err != nil {
		return 0, err
	}

	return int64(n), nil
}

// Name of the flag
func (f *FlagAccessor) Name() string {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("Name").String()
}

// Path of the flag
func (f *FlagAccessor) Path() string {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("Path").String()
}

// Usage of the flag
func (f *FlagAccessor) Usage() string {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("Usage").String()
}

// Value of the flag
func (f *FlagAccessor) Value() interface{} {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("Value").Interface()
}

// EnvVar of the flag
func (f *FlagAccessor) EnvVar() string {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("EnvVar").String()
}

// Hidden of the flag
func (f *FlagAccessor) Hidden() bool {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("Hidden").Bool()
}

// Validate validates the flag
func (f *FlagAccessor) Validate(ctx *Context) error {
	// FlagValidator validates a given flag
	type FlagValidator interface {
		Validate(ctx *Context) error
	}

	if validator, ok := f.Flag.(FlagValidator); ok {
		return validator.Validate(ctx)
	}

	return nil
}

// IsBoolFlag returns true if the flag is bool
func (f *FlagAccessor) IsBoolFlag() bool {
	// BoolFlag represents a boolean flag
	type BoolFlag interface {
		IsBoolFlag() bool
	}

	if flag, ok := f.Flag.(BoolFlag); ok {
		return flag.IsBoolFlag()
	}

	return false
}

// IsPathFlag returns true if the flag is path
func (f *FlagAccessor) IsPathFlag() bool {
	// PathFlag represents a boolean flag
	type PathFlag interface {
		IsPathFlag() bool
	}

	if flag, ok := f.Flag.(PathFlag); ok {
		return flag.IsPathFlag()
	}

	return false
}

// Reset resets the value
func (f *FlagAccessor) Reset() error {
	// FlagResetter resets a given flag
	type FlagResetter interface {
		Reset() error
	}

	if flag, ok := f.Flag.(FlagResetter); ok {
		return flag.Reset()
	}

	return nil
}

// FlagsByName is a slice of Flag
type FlagsByName []Flag

// Len returns the length of the slice
func (f FlagsByName) Len() int {
	return len(f)
}

// Less returns true if item at index i < item at index j
func (f FlagsByName) Less(i, j int) bool {
	return less(NewFlagAccessor(f[i]).Name(), NewFlagAccessor(f[j]).Name())
}

// Swap swaps two items
func (f FlagsByName) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

//go:generate counterfeiter -fake-name Validator -o ./fake/validator.go . Validator

// Validator converts values
type Validator interface {
	// Validate validates the value
	Validate(ctx *Context, value interface{}) error
}

var _ Validator = ValidatorFunc(nil)

// ValidatorFunc validates a flag
type ValidatorFunc func(ctx *Context, value interface{}) error

// Validate validates the value
func (fn ValidatorFunc) Validate(ctx *Context, value interface{}) error {
	return fn(ctx, value)
}
