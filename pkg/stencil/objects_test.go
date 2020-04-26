package stencil_test

import (
	"testing"

	"github.com/argots/stencil/pkg/stencil"
)

func TestObjectsSerialization(t *testing.T) {
	fs := fakeFS{
		files: map[string]string{".stencil/objects.json": "{}"},
	}
	s := stencil.New(discardLogger{}, discardLogger{}, nil, fs)
	if err := s.LoadObjects(); err != nil {
		t.Error("LoadObjects", err)
	}
	if err := s.SaveObjects(); err != nil {
		t.Error("SaveObjects", err)
	}
}
