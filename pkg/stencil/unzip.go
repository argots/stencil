package stencil

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"strings"
)

// Unzip extracts a single file from a zip source
func Unzip(src io.Reader, file string) io.ReadCloser {
	// unzip requires having the whole data in memory :(
	var buf bytes.Buffer
	size, err := io.Copy(&buf, src)
	if err != nil {
		return errReader{err}
	}

	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), size)
	if err != nil {
		return errReader{err}
	}
	for _, f := range r.File {
		if strings.EqualFold(file, f.Name) {
			result, err := f.Open()
			if err != nil {
				return errReader{err}
			}
			return result
		}
	}
	return errReader{errors.New("no such file: " + file)}
}
