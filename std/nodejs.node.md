# Install Node

This stencil installs Node into `./bin/node.`  This does not install `npm`.

## Version

The version of nodejs to use is defined by std.NodeVersion.

```go-template

{{ stencil.DefineString "std.NodeVersion" "Version of node (default v12.16.2)" }}
{{ $ver := (or (stencil.VarString "std.NodeVersion") "v12.16.2" ) }}

```

## Normalize for current OS and Arch

NodeJS uses a convention that deviates from stencil.OS and stencil.Arch:

```go-template

{{ $os := stencil.OS }}
{{ $arch := stencil.Arch }}

{{if (eq stencil.OS "windows") }}
  {{ $os = "win" }}
{{end}}
{{if (eq stencil.Arch "amd64") }}
  {{ $arch = "x64" }}
{{else if (eq stencil.Arch "386") }}
  {{ $arch = "x86" }}
{{end}}

```

## Fetch

```go-template

{{ $url := printf "https://nodejs.org/download/release/%s/node-%s-%s-%s.tar.gz" $ver $ver $os $arch }}
{{ $bin := printf "node-%s-%s-%s/bin/node" $ver $os $arch }}
{{ stencil.CopyFromArchive "nodejs" "./bin/node"  $url $bin }}

```
