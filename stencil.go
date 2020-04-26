package main

import (
	"flag"
	"log"
	"os"

	"github.com/argots/stencil/pkg/stencil"
)

func main() {
	baseDir, err := stencil.BaseDir()
	if err != nil {
		log.Fatal("stencil: basedir", err)
	}

	flags := flag.NewFlagSet("stencil", flag.ExitOnError)
	verbose := log.New(os.Stdout, "stencil: ", 0)
	errorl := log.New(os.Stderr, "stencil: ", 0)
	fs := &stencil.FS{BaseDir: baseDir, Verbose: verbose, Errorl: errorl}
	p := &stencil.ConsolePrompt{Stdin: os.Stdin, Stdout: os.Stdout}

	if err := stencil.New(verbose, errorl, p, fs).Main(flags, os.Args); err != nil {
		os.Exit(1)
	}
}
