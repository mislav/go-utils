package cli

import "strings"

// Flag is an option you can pass to the CLI via a long version, starting with --
// and optional a short one with only on -. It further has a help text and a type.
type Flag struct {
	Short string
	Long  string
	Help  string
	Ftype interface{}
}

// Parameter is an evaluated Flag
type Parameter struct {
	values   []string
	provided bool
	Flag
}

// Parameters provides methods to access all Parameter for a command
type Parameters struct {
	parameters map[string]*Parameter
}

// AddValue adds a value
func (f *Parameter) AddValue(v string) {
	f.values = append(f.values, v)
}

// IsProvided checks if the Parameter was provided when calling the CLI
func (f *Parameter) IsProvided() bool {
	return f.provided
}

// String returns the last provided value for this Parameter. If the Parameter
// was not provided it returns an empty string
func (f *Parameter) String() string {
	num := len(f.values)
	if num > 0 {
		return f.values[num-1]
	}
	return ""
}

// Bool returns the provided boolean values. If the Flag was set without value
// this return true, if the Flag was not set it returns false
func (f *Parameter) Bool() bool {
	val := strings.ToLower(f.String())
	return val == "true" || val == "t" || val == "1"
}

// IsProvided checks if a Parameter with the given name was provided when calling the CLI
func (f *Parameters) IsProvided(long string) bool {
	parameter := f.parameters[long]
	if parameter != nil {
		return parameter.IsProvided()
	}
	return false
}

// String returns the provided value for a Parameter. If the Parameter was not
// provided it returns the given fallback
func (f *Parameters) String(long, fallback string) string {
	parameter := f.parameters[long]
	if parameter != nil && parameter.IsProvided() {
		return parameter.String()
	}
	return fallback
}

// Bool returns the provided boolean values for a Parameter. If the Flag was
// set without value this return true, if the Flag was not set it returns false
func (f *Parameters) Bool(long string) bool {
	parameter := f.parameters[long]
	if parameter != nil {
		return parameter.Bool()
	}
	return false
}

// AddParameter adds a Parameter
func (f *Parameters) AddParameter(parameter *Parameter) {
	if f.parameters == nil {
		f.parameters = make(map[string]*Parameter)
	}
	f.parameters[parameter.Long] = parameter
}
