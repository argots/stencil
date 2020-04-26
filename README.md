# stencil

[![Test](https://github.com/argots/stencil/workflows/Test/badge.svg)](https://github.com/argots/stencil/actions?query=workflow%3ATest)
[![Lint](https://github.com/argots/stencil/workflows/Lint/badge.svg)](https://github.com/argots/stencil/actions?query=workflow%3ALint)
[![Go Report Card](https://goreportcard.com/badge/github.com/argots/stencil)](https://goreportcard.com/report/github.com/argots/stencil)

A simple package manager for static files and templates.

## Why another package manager?

Stencil is focused on two things -- allowing packages to cleanly
specify folder structures and allowing the files to be customized via
simple templates.

Stencil has a few goals that are different from other package
managers:

1. Only modify folders within the current directory.
2. Allow package writers to customize the files based on inputs chosen
the user (such as do you want "https" support?).
3. Allow users to modify the generated files without affecting the
ability to upgrade stencil packages -- they should properly merge much
like git itself.
4. Do not rely on specific environmental assets (such as the existence
of NodeJS or Ruby).  Any tools needed should be downloaded locally
(using github releases initially).

These design choices lead to a lot of duplication between stencil
projects, particularly with respect to tools. This is an explicit
tradeoff: the cost of duplication on disk is worth the benefit of
isolation between projects and environments.

## Install

The easiest way to install stencil is by downloading a release:

```sh
VERSION=v1.0.1
OS=darwin
ARCH=amd64
INSTALL_DIR=./bin
curl https://github.com/argots/stencil/releases/download/$VERSION/$OS_$ARCH.zip | tar -C $INSTALL_DIR -xvf - 
```

Note: on windows, this will install the binary without a .EXE
extension.  To get the right extension, do this instead:

```sh
curl https://github.com/argots/stencil/releases/download/$VERSION/$OS_$ARCH.zip | tar -C $INSTALL_DIR -xvf - stencil --transform "s/stencil/stencil.exe"
```

## Using stencil to install tools

Stencil follows the philosophy that all tools needed for a project
should be local to the project directory so that projects can be
isolated from each other.  So, installing tools would require that
binary tool releases are copied to a local directory, preferably
`./bin`.

For example, `nodejs` can be installed using the following script:

```sh
$ stencil pull git:git@github.com/argots/stencil#master/std/nodejs.node.stencil
```

This should copy the nodejs binary into `./bin`

## Todo

- [X] `stencil pull file.tpl` should run the file as a go template.
- [X] Add template function `stencil.CopyFile` to copy github URLs locally.
- [X] Add template variables `stencil.DefineBool("name", "prompt")` and `stencil.VarBool("name")`.
- [X] Add ability to modify variables `stencil -var XYZ=true`.
- [X] Add github releases.
- [X] Add support for downloading  github release via `stencil.CopyFromArchive`
- [ ] Add string variables and update nodejs install to ask for version.
- [ ] Add `stencil.ExtractArchive`
- [ ] Garbage collect unused files.
- [ ] Add unpull support.
- [ ] Add 3-way merge if git pull brings newer file and local file also modified.
- [ ] `stencil pull github.com/...` should fetch from public github (standard github url)
- [ ] Update `stencil.CopyFile(..)` to support relative github URLs
- [ ] `stencil pull github.com/...` should fetch from private github using ssh settings
- [ ] Add template variables `string("name", "prompt")` and `var("name")`
- [ ] `stencil pull` should pull latest versions of all URLs that were already pulled.
- [ ] Add nested templates support: `import(otherFile)`
- [ ] Add nested pull support `pull(args)`
- [ ] Add ability to use keyrings for secrets `secret("name", "prompt")` and `var("name")`
- [ ] Add ability to work with file patches inserted using markers
- [ ] Deal with diamond dependencies?
- [ ] Unsafe shell exec?