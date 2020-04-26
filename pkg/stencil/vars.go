package stencil

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
)

// NewVars creaate a new list of vars
func NewVars(f *flag.FlagSet, stdin io.Reader, stdout io.Writer) *Vars {
	var defs defsValue
	f.Var(&defs, "var", "bool_name or bool_name=yes/no/true/false")
	return &Vars{
		defs:       &defs,
		stdin:      stdin,
		stdout:     stdout,
		BoolDefs:   map[string]string{},
		StringDefs: map[string]string{},
		Bools:      map[string]bool{},
		Strings:    map[string]string{},
	}
}

// Vars holds named values.
type Vars struct {
	defs       *defsValue
	stdin      io.Reader
	stdout     io.Writer
	BoolDefs   map[string]string
	StringDefs map[string]string
	Bools      map[string]bool
	Strings    map[string]string
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

	val, err := v.readBool(prompt)
	if err != nil {
		return false, err
	}
	v.Bools[name] = val
	return val, nil
}

func (v *Vars) readBool(prompt string) (bool, error) {
	for {
		var s string
		fmt.Fprintf(v.stdout, "%s (Yes/No/True/False)? ", prompt)
		if _, err := fmt.Fscanln(v.stdin, &s); err != nil {
			return false, err
		}

		switch strings.ToLower(s) {
		case "t", "true", "y", "yes":
			return true, nil
		case "f", "false", "n", "no":
			return false, nil
		}
	}
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

	val, err := v.readString(prompt)
	if err != nil {
		return "", err
	}
	v.Strings[name] = val
	return val, nil
}

func (v *Vars) readString(prompt string) (string, error) {
	fmt.Fprintf(v.stdout, "%s? ", prompt)
	scanner := bufio.NewScanner(v.stdin)
	scanner.Scan()
	return scanner.Text(), scanner.Err()
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
