# Install GolangCI Linter

This stencil installs the golangci-lint into ./bin/golangici-lint.

## Usage

```bash
stencil pull git:git@github.com:argots/stencil.git/std/golangci-lint.md
```

## Version

The version of golagnci-lint is defined by std.GolangCILintVersion.

```go-template

{{ stencil.DefineString "std.GolangCILintVersion" "Version of golangci-lint (default v1.25.0)" }}
{{ $ver := (or (stencil.VarString "std.GolangCILintVersion") "v1.25.0") }}

```

## Normalize prefixes

Note that golangci-lint uses the version tag without the `v` within
the download url.  So we strip out the first character.

```go-template

{{ $ver = (slice $ver 1)}}

```

## Install

```go-template

{{ $os := stencil.OS }}
{{ $arch := stencil.Arch }}
{{ $url := printf "https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-%s-%s.tar.gz"  $ver $ver $os $arch }}
{{ $bin := printf "golangci-lint-%s-%s-%s/golangci-lint" $ver $os $arch }}
{{ stencil.CopyFromArchive "golangci-lint" "./bin/golangci-lint"  $url $bin }}

```
