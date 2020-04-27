# stencil

[![Test](https://github.com/argots/stencil/workflows/Test/badge.svg)](https://github.com/argots/stencil/actions?query=workflow%3ATest)
[![Lint](https://github.com/argots/stencil/workflows/Lint/badge.svg)](https://github.com/argots/stencil/actions?query=workflow%3ALint)
[![Go Report Card](https://goreportcard.com/badge/github.com/argots/stencil)](https://goreportcard.com/report/github.com/argots/stencil)

Stencil is a simple package manager for static files and templates.

## Contents
1. [Why another package manager?](#why-another-package-manager)
2. [Install](#install)
3. [Platforms supported](#platforms-supported)
4. [Stencil recipes are markdown and Go templating](#stencil-recipes-are-markdown-and-go-templating)
5. [Example templating with stencil](#example-templating-with-stencil)
6. [Code generation from markdowns](#code-generation-from-markdowns)
7. [Stencil variables](#stencil-variables)
8. [Status](#status)
9. [Todo](#todo)

## Why another package manager?

Stencil primarily came about from the pain of working with many
repositories.

Each repository typically has several needs:

- Setup CI/CD
- Setup linter rules, Makefile
- Setup code scaffodling (such as
[create-react-app](https://github.com/facebook/create-react-app)). 
- Setup standard tools with specific versions needed without
conflicting with tools in other repositories.
- Ability to update the templates for these and have them apply to all
derived repositories.

A typical approach is to use a [github template
repository](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/creating-a-template-repository)
but this only solves a small set of the problems identified above.

Stencil sets out to solve these problems:

1. Provide a rich declarative templating system that allows creating a
local disk layout for a *workspace*.
2. Provide a [mad-lib](https://en.wikipedia.org/wiki/Mad_Libs) style
interactive configurability by prompting users for choices and
remembering those choices for future runs.
3. Provide the ability to change templates and pull these changes
locally *a la Git* with proper three-way merges.
4. Allow multiple templates to be used within a *workspace.*


Stencil makes a few opinionated choices:

1. **Strong Isolation:**  Stencil recipes cannot modify any files or
the environment outside a workspace.  This includes any tools
(including language environments).  This has the side-effect that
tool-chains are duplicated but the expliclit tradeoff is we getclean
isolation and avoid the complexity and unpredictability of
shared tools.  At some point, Stencil may get smart enough to use
[Content addressable
storage](https://en.wikipedia.org/wiki/Content-addressable_storage) techniques.
2. **Environmental idependence:** Stencil recipes cannot execute
arbitrary programs -- they can only modify local files based on user
choices and can only download executables.  This is a severe
limitation in theory but most tools provide downloadable binaries in
practice.  Another approach here is to setup tool repositories that
hold the release binaries (such as using [github
releases](https://help.github.com/en/enterprise/2.13/user/articles/creating-releases)).
3. **Declarative Setup:** Stencil aims to adopt a somewhat [React-like
reconciliaton](https://reactjs.org/docs/reconciliation.html) process.
All changes done by Stencil are remembered between runs and so any
files that are not written by a specific run are automatically
cleaned up. This incurs a performance cost but makes it easier to
write Stencil recipes and maintain repositories.  The actual APIs do
not have a declarative feel beyond requiring a "key" argument to make
it easier to use.
4. **Conflict Resolution:** Stencil aims to handle upstream changes of
templates gracefully even when the local files are modified by
invoking a three-way merge.  Much like Git, if this fails, the
local files are left in a conflicted state (with manual fixups as
needed).

## Install

The easiest way to install stencil is by downloading a release:

```sh
VERSION=v1.0.1
OS=darwin
ARCH=amd64
INSTALL_DIR=./bin
curl https://github.com/argots/stencil/releases/download/$VERSION/$OS_$ARCH.zip | tar -C $INSTALL_DIR -xvf - 
```

## Platforms supported

Stencil is being developed with cross-platform support in mind but
this is is not a priority yet.  MacOS and Linux are the initial
priorities.

For instance, even installing stencil on windows requires a slight
change to get the binary to have a `.EXE` extension: 

```sh
curl https://github.com/argots/stencil/releases/download/$VERSION/$OS_$ARCH.zip | tar -C $INSTALL_DIR -xvf - stencil --transform "s/stencil/stencil.exe"
```

There are no doubt other such minor changes needed.

## Stencil recipes are markdown and Go templating

Please see
[golangci-lint](https://github.com/argots/stencil/blob/master/std/golangci-lint.md)
or
[NodeJS](https://github.com/argots/stencil/blob/master/std/nodejs.node.md)
or [Go](https://github.com/argots/stencil/blob/master/std/golang.md)
for example recipes.

These recipes are all markdown files to promote readability.  Stencil
uses the Go templating engine underneath even for the recipes and it
discards the output file.  The actual files that are copied are
specified within the markdown in the code fences where the
Stencil-provided functions are called to accomplish this.  This
provides a good scripting environment for stencil.

These templates are invoked with the following syntax (which is
slightly different from the github URL associated with these):

```bash
stencil pull git:git@github.com:argots/stencil.git/std/nodejs.node.md
```

## Example templating with stencil

This example uses local files to illustrate but this can also work
with git urls.

First we create a recipe that prompts the user for a package name and
then creates `./pkg/<name>/<name>.go` file:

```bash
cat > my_recipe.md <<EOF
This copies my_file.go to ./pkg/<name>/<name>.go.

Lets first ask the user for what <name> should be:

{{ stencil.DefineString "pkg.Name" "Whats the pkg name? (default boo)" }}
{{ \$name := stencil.VarString "pkg.Name" }}
{{ \$dir := printf "./pkg/%s/%s.go" \$name \$name }}
{{ stencil.CopyFile "my_recipe" \$dir "my_file.go.tpl" }}
EOF
```

The recipe references a `my_file.go.tpl`, so lets create this template:

```bash
$ cat > my_file.go.tpl <<EOF
package {{ stencil.VarString "pkg.Name" }}

import "fmt"

func init() {
	fmt.Println("Yay")
}

EOF
```

Now we can test it all:


```bash
$ stencil pull my_recipe.md
```

## Code generation from markdowns

Stencil also supports a `stencil.CopyMarkdownSnippet` function which
allows code to be embedded within the recipe itself with code fences.

This is shown in the [snippets.md
example](https://github.com/argots/stencil/blob/master/examples/snippets.md).

This recipe includes a code fence tagged with `golang` which is then
used in the `stencil.CopyMarkdownSnippets` instruction -- the pattern
provided to this function is used to filter out the relevant code
fences and only include those in the `./example.go` file genrated.

The following command invokes the recipe above:

```bash
stencil pull git:git@github.com:argots/stencil.git/examples/snippets.md
```

## Stencil variables

Stencil variables are meant to hold configurable things like version
of Node or package name etc.  These are prompted for only once with
subsequent attempts reusing the last value (which is saved in the file
`.stencil/objects.json`).

Variables can be forcibly changed by passing in `--var Name` (for
booleans) and `--var Name=Value` (for all types).

Within a template, variables must first be defined using
`stencil.DefineString "name" "prompt"` or `stencil.DefineBool "name"
"prompt"` before the current value is fetched.  The current value can
be fetched using `stencil.VarString "name"` or `stencil.VarBool
"name"`.

The variable name can only be discovered by looking at the recipe or
looking at `.stencil/objects.json` (which has a Strings and Bools key
with the associated values).

## Status

This is still unstable.  In particular, the APIs may change slightly
as they do not compose very well still.

## Todo

- [X] `stencil pull file.tpl` should run the file as a go template.
- [X] Add template function `stencil.CopyFile` to copy github URLs locally.
- [X] Add template variables `stencil.DefineBool("name", "prompt")` and `stencil.VarBool("name")`.
- [X] Add template variables `stencil.DefineString` and `stencil.VarString`
- [X] Add ability to modify variables `stencil -var XYZ=true`.
- [X] Add github releases.
- [X] Add support for downloading  github release via `stencil.CopyFromArchive`
- [X] Add `stencil.CopyManyFromArchive`
- [X] Add string variables and update nodejs install to ask for version.
- [X] Garbage collect unused files.
- [X] `stencil pull git:git@github.com:repo/user/path` should fetch from public github (standard github url)
- [X] `stencil pull ...` should fetch from private github using ssh settings
- [X] Add `stencil.CopyMarkdownSnippets` support
- [ ] Add unpull support.
- [ ] `stencil pull` should pull latest versions of all URLs that were already pulled.
- [ ] Add 3-way merge if git pull brings newer file and local file also modified.
- [ ] Add ability to look at all variable values.
- [ ] Add nested templates support: `import(otherFile)`
- [ ] Update `stencil.CopyFile` to support relative github URLs
- [ ] Add nested pull support `pull(args)`
- [ ] Add ability to use keyrings for secrets
- [ ] Add ability to work with file patches inserted using markers
- [ ] Deal with diamond dependencies?
- [ ] Unsafe shell exec?
- [ ] Other template engines than the default Go? Other script languages?
