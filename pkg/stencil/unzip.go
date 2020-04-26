package stencil

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
)

// Unzip visits all files in a zip archive.
func Unzip(src io.Reader, visit func(string, func() io.ReadCloser) error) error {
	// unzip requires having the whole data in memory :(
	var buf bytes.Buffer
	size, err := io.Copy(&buf, src)
	if err != nil {
		return err
	}

	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), size)
	if err != nil {
		return err
	}
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "/") {
			continue
		}

		if err := visit(f.Name, unzipOpener(f)); err != nil {
			return err
		}
	}
	return nil
}

func unzipOpener(f *zip.File) func() io.ReadCloser {
	return func() io.ReadCloser {
		r, err := f.Open()
		if err != nil {
			return errReader{err}
		}
		return r
	}
}
