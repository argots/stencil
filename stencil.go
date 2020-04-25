package main

import (
	"flag"
	"log"
	"os"

	"github.com/argots/stencil/pkg/stencil"
)

func main() {
	f := flag.NewFlagSet("stencil", flag.ExitOnError)

	verbose := log.New(os.Stdout, "stencil: ", 0)
	errorl := log.New(os.Stderr, "stencil: ", 0)
	cache, err := stencil.NewCache(f, os.Stdin, os.Stdout)
	if err != nil {
		panic(err)
	}

	fs := &stencil.FS{Verbose: verbose, Errorl: errorl}

	if err := stencil.New(verbose, errorl, cache, fs).Main(f, os.Args); err != nil {
		os.Exit(1)
	}
}
