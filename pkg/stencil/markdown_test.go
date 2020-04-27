package stencil_test

import (
	"strings"
	"testing"

	"github.com/argots/stencil/pkg/stencil"
)

func TestMarkdownFilter(t *testing.T) {
	md := `
This is a sample markdown.

fenceGOLANG
stuff within a f.nce
fence

fenceGOLANG
stuff within a f.nce
fence
`
	md = strings.ReplaceAll(md, "fence", "```")
	m := &stencil.Markdown{}
	result, err := m.FilterMarkdown(md, "G.L")
	if err != nil {
		t.Fatal("err", err)
	}
	if result != "stuff within a f.nce\nstuff within a f.nce\n" {
		t.Error("Got unexpected result", result)
	}
}
