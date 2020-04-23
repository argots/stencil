package stencil

import (
	"io/ioutil"
)

// FS implements a FileSystem interface.
type FS struct {
}


// Read reads the contents of the path and returns them as bytes.
func (fs *FS) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
