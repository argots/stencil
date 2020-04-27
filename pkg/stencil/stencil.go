// Package stencil implements a template package manager.
package stencil

import (
	"bytes"
	"errors"
	"flag"
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
	Remove(path string) error
	RemoveAll(path string) error
}

// Prompter is the generic interface to prompt and fetch info
// interactively.
type Prompter interface {
	PromptBool(prompt string) (bool, error)
	PromptString(prompt string) (string, error)
}

// New creates a new stencil manager.
func New(verbose, errorl Logger, p Prompter, fs FileSystem) *Stencil {
	s := &Stencil{
		State:  map[string]interface{}{},
		Funcs:  map[string]interface{}{},
		Printf: verbose.Printf,
		Errorf: func(fmt string, v ...interface{}) error {
			errorl.Printf(fmt, v...)
			err, _ := v[len(v)-1].(error)
			return err
		},
		FileSystem: fs,
		Prompter:   p,
		Binary:     Binary{},
		Objects: Objects{
			Before:       &Objects{},
			Pulls:        map[string]bool{},
			Files:        map[string]*FileObj{},
			FileArchives: map[string]*FileArchiveObj{},
			Bools:        map[string]bool{},
			Strings:      map[string]string{},
		},
		Vars: Vars{
			BoolDefs:   map[string]string{},
			StringDefs: map[string]string{},
		},
		Markdown: Markdown{},
	}
	s.Binary.Stencil = s
	s.Objects.Stencil = s
	s.Vars.Stencil = s
	s.Markdown.Stencil = s
	s.Funcs["stencil"] = func() interface{} {
		return s
	}
	return s
}

// Stencil maintains all the state for managing a single directory.
type Stencil struct {
	State  map[string]interface{}
	Funcs  map[string]interface{}
	Printf func(format string, v ...interface{})
	Errorf func(format string, v ...interface{}) error
	FileSystem
	Prompter
	Env
	Binary
	Objects
	Vars
	Markdown
}

// Main implements the main program.
func (s *Stencil) Main(f *flag.FlagSet, args []string) error {
	errMissingArg := errors.New("missing arg")
	f.Usage = func() {
		s.Printf(`Usage: stencil [options] commands
  commands:
    pull url_or_file -- add url to pulls and sync
    rm url_or_fil    -- remove url from pulls and sync
    sync             -- update all existing pulls
`)
		f.PrintDefaults()
	}
	s.Vars.Init(f)
	if err := f.Parse(args[1:]); err != nil {
		return s.Errorf("flagset parse", err)
	}

	switch f.Arg(0) {
	case "pull":
		if f.Arg(1) != "" {
			s.Printf("Adding %s\n", f.Arg(1))
			return s.run(f.Arg(1), "")
		}
		return s.Errorf("pull requires a url or path to a recipe %v\n", errMissingArg)
	case "sync":
		s.Printf("Updating all pulled recipes\n")
		return s.run("", "")
	case "rm":
		if f.Arg(1) != "" {
			s.Printf("Removing %s\n", f.Arg(1))
			return s.run("", f.Arg(1))
		}
		return s.Errorf("rm requires a url or path to a recipe %v\n", errMissingArg)
	case "":
		f.Usage()
		return nil
	}
	return s.Errorf("%v", errors.New("unknown command: "+f.Arg(0)))
}

func (s *Stencil) run(add, rm string) error {
	if err := s.Objects.LoadObjects(); err != nil {
		return s.Errorf("LoadObjects %v\n", err)
	}
	for pull := range s.Before.Pulls {
		if pull != rm {
			s.Objects.addPull(pull)
		}
	}
	if add != "" {
		s.Objects.addPull(add)
	}
	for pull := range s.Pulls {
		s.Printf("Pulling %s\n", pull)
		if err := s.Run(pull); err != nil {
			return s.Errorf("Run: %v\n", err)
		}
	}
	if err := s.GC(); err != nil {
		return s.Errorf("GC %v\n", err)
	}

	return s.SaveObjects()
}

// CopyFile copies a url to a local file.
func (s *Stencil) CopyFile(key, localPath, url string) error {
	s.Printf("copying %s to %s, key (%s)\n", url, localPath, key)
	s.Objects.addFile(key, localPath, url)

	data, err := s.Execute(url)
	if err != nil {
		return s.Errorf("Error reading %s %v\n", url, err)
	}
	return s.Write(localPath, []byte(data), 0666)
}

// Run runs a template discarding the output.
func (s *Stencil) Run(source string) error {
	_, err := s.Import(source)
	return err
}

// Import imports a template after applying it.
func (s *Stencil) Import(source string) (string, error) {
	s.Printf("Running import %s\n", source)
	return s.Execute(source)
}

// Execute is like Run but it returns the output.
func (s *Stencil) Execute(source string) (string, error) {
	return s.executeFilter(source, nil)
}

func (s *Stencil) executeFilter(source string, filter func(string) (string, error)) (string, error) {
	s.Printf("Executing %s\n", source)

	data, err := s.Read(source)
	if err != nil {
		return "", s.Errorf("Error reading %s: %v\n", source, err)
	}

	if filter != nil {
		str, err := filter(string(data))
		if err != nil {
			return "", s.Errorf("Error filtering %s: %v", source, err)
		}
		data = []byte(str)
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
