package stencil

import "runtime"

// Env implements some standard environent accessors for stencil.
type Env struct{}

// OS returns runtime.GOOS.
func (e Env) OS() string {
	return runtime.GOOS
}

// Arch returns runtime.GOARCH
func (e Env) Arch() string {
	return runtime.GOARCH
}
