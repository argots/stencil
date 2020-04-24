package main

import (
	"log"
	"os"

	"github.com/argots/stencil/pkg/stencil"
)

func main() {
	verbose := log.New(os.Stdout, "stencil: ", 0)
	errorl := log.New(os.Stderr, "stencil: ", 0)
	fs := &stencil.FS{Verbose: verbose, Errorl: errorl}

	if err := stencil.New(verbose, errorl, fs).Main(); err != nil {
		os.Exit(1)
	}
}
