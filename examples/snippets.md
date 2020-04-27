# Markdown With Embedded Code Snippets

This stencil example illustrates a recipe which includes code snippets.

The following snippet has been taggged as golang.

```golang
package example

import "fmt"

func init() {
	fmt.Println("This is an example")
}
```

Now, the following stencil code will effectively copy this snippet and
write it to a file named `./example.go`:

```go-template
{{ $src := "git:git@github.com:argots/stencil.git/examples/snippets.md" }}
{{ $pattern := "golang" }}
{{ stencil.CopyMarkdownSnippets "ex" "./example.go" $src $pattern }}
```
