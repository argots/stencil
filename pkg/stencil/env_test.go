package stencil_test

import (
	"runtime"
	"testing"

	"github.com/argots/stencil/pkg/stencil"
)

func TestEnv(t *testing.T) {
	var env stencil.Env
	if env.OS() != runtime.GOOS {
		t.Error("Env.OS", env.OS())
	}
	if env.Arch() != runtime.GOARCH {
		t.Error("Env.Arch", env.Arch())
	}
}
