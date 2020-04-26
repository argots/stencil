package stencil

import "io"

var _ io.ReadCloser = errReader{}

type errReader struct {
	err error
}

func (e errReader) Read(data []byte) (int, error) {
	return 0, e.err
}

func (e errReader) Close() error {
	return nil
}
