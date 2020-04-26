// Package stencil implements a template package manager.
package stencil

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"text/template"
)

// Logger is the interface used by stencil to log messages.
//
// Use log.New(...) to create an appropriate logger.
type Logger interface {
	Printf(format string, v ...interface{})
}

// FileSystem is the generic file system used by stencil.
// Use FS{} or a custom implementation.
type FileSystem interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte, mode os.FileMode) error
}

// New creates a new stencil manager.
func New(verbose, errorl Logger, c *Cache, fs FileSystem) *Stencil {
	s := &Stencil{
		State:  map[string]interface{}{},
		Funcs:  map[string]interface{}{},
		Cache:  c,
		Printf: verbose.Printf,
		Errorf: func(fmt string, v ...interface{}) error {
			errorl.Printf(fmt, v...)
			err, _ := v[len(v)-1].(error)
			return err
		},
		FileSystem: fs,
		Binary:     Binary{fs},
	}
	s.Funcs["stencil"] = func() interface{} {
		return s
	}
	return s
}

// Stencil maintains all the state for managing a single directory.
type Stencil struct {
	State map[string]interface{}
	Funcs map[string]interface{}
	*Cache
	Printf func(format string, v ...interface{})
	Errorf func(format string, v ...interface{}) error
	FileSystem
	Env
	Binary
}

// Main implements the main program
func (s *Stencil) Main(f *flag.FlagSet, args []string) error {
	f.Usage = func() {
		fmt.Fprint(s.Cache.Vars.stdout, "Usage: stencil [options] pull github-url")
		f.PrintDefaults()
	}
	f.SetOutput(s.Cache.Vars.stdout)
	if err := f.Parse(args[1:]); err != nil {
		return s.Errorf("flagset parse", err)
	}

	switch f.Arg(0) {
	case "pull":
		return s.Run(f.Arg(1))
	case "":
		f.Usage()
		return nil
	}
	return s.Errorf("%v", errors.New("unknown command: "+f.Arg(0)))
}

// CopyFile copies a url to a local file.
func (s *Stencil) CopyFile(key, localPath, url string) error {
	s.Printf("copying %s to %s, key (%s)\n", url, localPath, key)
	data, err := s.Read(url)
	if err != nil {
		return s.Errorf("Error reading %s %v\n", url, err)
	}
	return s.Write(localPath, data, 0666)
}

// Run runs a template discarding the output.
func (s *Stencil) Run(source string) error {
	_, err := s.Import(source)
	return err
}

// Import imports a template after applying it.
func (s *Stencil) Import(source string) (string, error) {
	s.Printf("Running import %s\n", source)

	data, err := s.Read(source)
	if err != nil {
		return "", s.Errorf("Error reading %s: %v\n", source, err)
	}

	t, err := template.New(source).Funcs(s.Funcs).Parse(string(data))
	if err != nil {
		return "", s.Errorf("Error parsing %s: %v\n", source, err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, s.State)
	if err != nil {
		return "", s.Errorf("Error executing %s: %v\n", source, err)
	}

	return buf.String(), nil
}
