package stencil

import (
	"archive/tar"
	"io"
)

// Untar visits all files in a tar archive.
func Untar(src io.Reader, visit func(string, func() io.ReadCloser) error) error {
	r := tar.NewReader(src)
	for {
		next, err := r.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		}
		if next.Typeflag != tar.TypeReg {
			continue
		}

		if err = visit(next.Name, (untarReadCloser{r}).self); err != nil {
			return err
		}
	}
}

type untarReadCloser struct {
	io.Reader
}

func (untarReadCloser) Close() error {
	return nil
}

func (u untarReadCloser) self() io.ReadCloser {
	return u
}
