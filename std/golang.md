# Install Go

This stencil installs a standard golang distribution.  Go requries the
whole zip file unfortunately, so this installs the actual go binary
under ./bin/go/bin/go

## Version

The version of Go to use is defined by the std.GoVersion variable.

```go-template

{{ stencil.DefineString "std.GoVersion" "Version of Go (default v1.14.2)" }}
{{ $ver := (or (stencil.VarString "std.GoVersion") "v1.14.2" ) }}

```

## Install

```go-template

{{ $os := stencil.OS }}
{{ $arch := stencil.Arch }}
{{ $url := printf "https://dl.google.com/go/%s.%s-%s.tar.gz" $ver $os $arch }}
{{ stencil.CopyManyFromArchive "golang" "./bin/"  $url "**" }}

```
