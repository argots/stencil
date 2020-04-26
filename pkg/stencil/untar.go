package stencil

import (
	"archive/tar"
	"errors"
	"io"
	"strings"
)

// Untar extracts a single file from a tar archive.
func Untar(src io.Reader, file string) io.Reader {
	r := tar.NewReader(src)
	for {
		next, err := r.Next()
		switch {
		case err == io.EOF:
			return errReader{errors.New("file not found " + file)}
		case err != nil:
			return errReader{err}
		case strings.EqualFold(file, next.Name):
			return r
		}
	}
}
