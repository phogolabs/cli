package cli

import "fmt"

// ValidationFn validates a flag
type ValidationFn func(f Flag) error

// FlagDefinition is the flag's definition
type FlagDefinition struct {
	Name     string
	Usage    string
	EnvVar   string
	FilePath string
	Metadata map[string]string
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
	Definition() *FlagDefinition
	Validate() error
}

var _ Flag = &StringFlag{}

// StringFlag is a flag with type string
type StringFlag struct {
	Name         string
	Usage        string
	EnvVar       string
	FilePath     string
	Value        string
	Metadata     map[string]string
	Hidden       bool
	Required     bool
	ValidationFn ValidationFn
}

// String returns the value as string
func (f *StringFlag) String() string {
	return f.Value
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

	if f.ValidationFn != nil {
		return f.ValidationFn(f)
	}

	return nil
}

// Definition returns the flag's definition
func (f *StringFlag) Definition() *FlagDefinition {
	return &FlagDefinition{
		Name:     f.Name,
		Usage:    f.Usage,
		EnvVar:   f.EnvVar,
		FilePath: f.FilePath,
		Metadata: f.Metadata,
	}
}

// RequiredErr returns the required error
func RequiredErr(flag string) error {
	return fmt.Errorf("cli: flag -%s is missing", flag)
}
