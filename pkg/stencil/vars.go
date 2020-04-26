package stencil

import (
	"errors"
	"flag"
	"strings"
)

// Vars holds named values.
type Vars struct {
	*Stencil
	defs       defsValue
	BoolDefs   map[string]string
	StringDefs map[string]string
}

// Init initializes vars.  Must be calleed for flag.Parse.
func (v *Vars) Init(f *flag.FlagSet) {
	f.Var(&v.defs, "var", "bool_name or bool_name=yes/no/true/false or string_name=value")
}

// DefineBool defines a boolean variable name.
func (v *Vars) DefineBool(name, prompt string) error {
	if _, ok := v.BoolDefs[name]; ok {
		return errors.New("redefiniton of " + name)
	}
	v.BoolDefs[name] = prompt
	return nil
}

// DefineString defines a string variable name.
func (v *Vars) DefineString(name, prompt string) error {
	if _, ok := v.StringDefs[name]; ok {
		return errors.New("redefiniton of " + name)
	}
	v.StringDefs[name] = prompt
	return nil
}

// VarBool fetches the value for the named boolean.  If the value is
// not present either via --var or via a previous invocation, the
// value is prompted for using the prompt in the definition.
// Any --var use overrides default values present from previous
// invocations.
func (v *Vars) VarBool(name string) (bool, error) {
	prompt, ok := v.BoolDefs[name]
	if !ok {
		return false, errors.New("undefined variable: " + name)
	}

	if val, ok := v.defs.bools[name]; ok {
		v.Bools[name] = val
		return val, nil
	}

	if val, ok := v.Bools[name]; ok {
		return val, nil
	}

	if val, ok := v.Before.Bools[name]; ok {
		v.Bools[name] = val
		return val, nil
	}

	val, err := v.PromptBool(prompt)
	if err == nil {
		v.Bools[name] = val
	}
	return val, err
}

// VarString fetches the value for the named variable.  If the value is
// not present either via --var or via a previous invocation, the
// value is prompted for using the prompt in the definition.
// Any --var use overrides default values present from previous
// invocations.
func (v *Vars) VarString(name string) (string, error) {
	prompt, ok := v.StringDefs[name]
	if !ok {
		return "", errors.New("undefined variable: " + name)
	}

	if val, ok := v.defs.strings[name]; ok {
		v.Strings[name] = val
		return val, nil
	}

	if val, ok := v.Strings[name]; ok {
		return val, nil
	}

	if val, ok := v.Before.Strings[name]; ok {
		v.Strings[name] = val
		return val, nil
	}

	val, err := v.PromptString(prompt)
	if err != nil {
		return "", err
	}

	v.Strings[name] = val
	return val, nil
}

type defsValue struct {
	bools   map[string]bool
	strings map[string]string
}

func (d *defsValue) String() string {
	return "serializings defs NYI"
}

func (d *defsValue) Set(value string) error {
	if d.bools == nil {
		d.bools = map[string]bool{}
		d.strings = map[string]string{}
	}

	parts := strings.Split(value, "=")

	if len(parts) == 1 {
		d.bools[value] = true
		return nil
	}

	switch strings.ToLower(parts[1]) {
	case "y", "yes", "t", "true":
		d.bools[parts[0]] = true
		return nil
	case "n", "no", "f", "false":
		d.bools[parts[0]] = false
		return nil
	}
	d.strings[parts[0]] = parts[1]
	return nil
}
