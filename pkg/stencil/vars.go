package stencil

import (
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
		defs:     &defs,
		stdin:    stdin,
		stdout:   stdout,
		BoolDefs: map[string]string{},
		Bools:    map[string]bool{},
	}
}

// Vars holds named values.
type Vars struct {
	defs     *defsValue
	stdin    io.Reader
	stdout   io.Writer
	BoolDefs map[string]string
	Bools    map[string]bool
}

func (v *Vars) DefineBool(name, prompt string) error {
	v.BoolDefs[name] = prompt
	return nil
}

func (v *Vars) VarBool(name string) (bool, error) {
	prompt, ok := v.BoolDefs[name]
	if !ok {
		return false, errors.New("undefined variable: " + name)
	}

	if val, ok := v.Bools[name]; ok {
		return val, nil
	}

	if val, ok := v.defs.bools[name]; ok {
		v.Bools[name] = val
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

type defsValue struct {
	bools map[string]bool
}

func (d *defsValue) String() string {
	return "serializings defs NYI"
}

func (d *defsValue) Set(value string) error {
	if d.bools == nil {
		d.bools = map[string]bool{}
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
	panic("NYI")
}
