package cli

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/oliveagle/jsonpath"
	"gopkg.in/yaml.v2"
)

//go:generate counterfeiter -fake-name Converter -o ./fake/converter.go . Converter

// Converter converts values
type Converter interface {
	// Convert converts the value
	Convert(interface{}) (interface{}, error)
}

var _ Converter = ConverterFunc(nil)

// ConverterFunc converts values
type ConverterFunc func(interface{}) (interface{}, error)

// Convert converts the value
func (fn ConverterFunc) Convert(v interface{}) (interface{}, error) {
	return fn(v)
}

var _ Converter = JSONPath("")

// JSONPath converts a value from JSON by Path
type JSONPath string

// Convert converts the value
func (path JSONPath) Convert(v interface{}) (interface{}, error) {
	return jsonpath.JsonPathLookup(v, string(path))
}

//go:generate counterfeiter -fake-name Validator -o ./fake/validator.go . Validator

// Validator converts values
type Validator interface {
	// Validate validates the value
	Validate(interface{}) error
}

var _ Validator = ValidatorFunc(nil)

// ValidatorFunc validates a flag
type ValidatorFunc func(interface{}) error

// Validate validates the value
func (fn ValidatorFunc) Validate(v interface{}) error {
	return fn(v)
}

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
	Validate() error
}

var _ Flag = &StringFlag{}

// StringFlag is a flag with type string
type StringFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     string
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *StringFlag) Validate() error {
	if f.Required {
		if f.Value == "" {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

var _ Flag = &StringSliceFlag{}

// StringSliceFlag is a flag with type *StringSlice
type StringSliceFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     []string
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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

// Reset resets the value
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
func (f *StringSliceFlag) Validate() error {
	if f.Required {
		if f.Value == nil || len(f.Value) == 0 {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

var _ Flag = &BoolFlag{}

// BoolFlag is a flag with type bool
type BoolFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     bool
	Metadata  map[string]string
	Hidden    bool
	Converter Converter
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
func (f *BoolFlag) Validate() error {
	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// URLFlag is a flag with type url.URL
type URLFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     *url.URL
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *URLFlag) Validate() error {
	if f.Required {
		if f.Value == nil {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// JSONFlag is a flag with type json document
type JSONFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     interface{}
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
	if f.Value == nil {
		f.Value = make(map[string]interface{})
	}

	return json.Unmarshal([]byte(value), &f.Value)
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *JSONFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *JSONFlag) Validate() error {
	if f.Required {
		if f.Value == nil {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// YAMLFlag is a flag with type yaml document
type YAMLFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     interface{}
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
	if f.Value == nil {
		f.Value = make(map[string]interface{})
	}

	return yaml.Unmarshal([]byte(value), &f.Value)
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *YAMLFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *YAMLFlag) Validate() error {
	if f.Required {
		if f.Value == nil {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// XMLFlag is a flag with type XMLDocument
type XMLFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     interface{}
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
	if f.Value == nil {
		return nil
	}
	return xml.Unmarshal([]byte(value), &f.Value)
}

// Get is a function that allows the contents of a Value to be retrieved.
// It wraps the Value interface, rather than being part of it, because it
// appeared after Go 1 and its compatibility rules. All Value types provided
// by this package satisfy the Getter interface.
func (f *XMLFlag) Get() interface{} {
	return f.Value
}

// Validate validates the flag
func (f *XMLFlag) Validate() error {
	if f.Required {
		if f.Value == nil {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// TimeFlag is a flag with type time.Time
type TimeFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Format    string
	Value     time.Time
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *TimeFlag) Validate() error {
	if f.Required {
		if f.Value.IsZero() {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// DurationFlag is a flag with type time.Duration
type DurationFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     time.Duration
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *DurationFlag) Validate() error {
	if f.Required {
		if f.Value == 0 {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// IntFlag is a flag with type int
type IntFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     int
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *IntFlag) Validate() error {
	if f.Required {
		if f.Value == 0 {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// Int64Flag is a flag with type int64
type Int64Flag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     int64
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *Int64Flag) Validate() error {
	if f.Required {
		if f.Value == 0 {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// UIntFlag is a flag with type uint64
type UIntFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     uint
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *UIntFlag) Validate() error {
	if f.Required {
		if f.Value == 0 {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// UInt64Flag is a flag with type uint
type UInt64Flag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     uint64
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *UInt64Flag) Validate() error {
	if f.Required {
		if f.Value == 0 {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// Float32Flag is a flag with type float32
type Float32Flag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     float32
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *Float32Flag) Validate() error {
	if f.Required {
		if f.Value == 0 {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// Float64Flag is a flag with type float64
type Float64Flag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     float64
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *Float64Flag) Validate() error {
	if f.Required {
		if f.Value == 0 {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// IPFlag is a flag with type net.IP
type IPFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     net.IP
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *IPFlag) Validate() error {
	if f.Required {
		if f.Value == nil {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// HardwareAddrFlag is a flag with type net.HardwareAddr
type HardwareAddrFlag struct {
	Name      string
	Usage     string
	EnvVar    string
	FilePath  string
	Value     net.HardwareAddr
	Metadata  map[string]string
	Hidden    bool
	Required  bool
	Converter Converter
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
func (f *HardwareAddrFlag) Validate() error {
	if f.Required {
		if f.Value == nil {
			return RequiredErr(f.Name)
		}
	}

	if f.Validator != nil {
		return f.Validator.Validate(f.Value)
	}

	return nil
}

// FlagAccessor access the flag's field
type FlagAccessor struct {
	Flag Flag
}

// Value of the flag
func (f *FlagAccessor) Value() interface{} {
	return f.Flag.Get()
}

// SetValue sets the value
func (f *FlagAccessor) SetValue(v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = f.error(r)
		}
	}()

	if converter := f.Converter(); converter != nil {
		if v, err = converter.Convert(v); err != nil {
			return err
		}
	}

	value := reflect.ValueOf(v)

	switch value.Kind() {
	case reflect.String:
		return f.Flag.Set(value.String())
	default:
		flag := reflect.ValueOf(f.Flag)
		flag = reflect.Indirect(flag)
		field := flag.FieldByName("Value")

		if field.CanSet() {
			field.Set(value)
		}
	}

	return err
}

// Converter returns the converter
func (f *FlagAccessor) Converter() Converter {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	field := value.FieldByName("Converter")

	if converter, ok := field.Interface().(Converter); ok {
		return converter
	}

	return nil
}

// Name of the flag
func (f *FlagAccessor) Name() string {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("Name").String()
}

// Usage of the flag
func (f *FlagAccessor) Usage() string {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("Usage").String()
}

// EnvVar of the flag
func (f *FlagAccessor) EnvVar() string {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("EnvVar").String()
}

// FilePath of the flag
func (f *FlagAccessor) FilePath() string {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("FilePath").String()
}

// Metadata of the flag
func (f *FlagAccessor) Metadata() map[string]string {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)

	metadata, ok := value.FieldByName("Metadata").Interface().(map[string]string)
	if !ok {
		return nil
	}
	return metadata
}

// MetaKey returns a metadata by key
func (f *FlagAccessor) MetaKey(path string) string {
	if value, ok := f.Metadata()[path]; ok {
		return value
	}
	return ""
}

// Hidden of the flag
func (f *FlagAccessor) Hidden() bool {
	value := reflect.ValueOf(f.Flag)
	value = reflect.Indirect(value)
	return value.FieldByName("Hidden").Bool()
}

// Reset the value
func (f *FlagAccessor) Reset() error {
	if r, ok := f.Flag.(resetter); ok {
		if err := r.Reset(); err != nil {
			return err
		}
	}

	return nil
}

func (f *FlagAccessor) error(v interface{}) error {
	switch err := v.(type) {
	case *reflect.ValueError:
		return err
	default:
		return fmt.Errorf("%v", err)
	}
}

// RequiredErr returns the required error
func RequiredErr(flag string) error {
	return fmt.Errorf("cli: flag -%s is missing", flag)
}
