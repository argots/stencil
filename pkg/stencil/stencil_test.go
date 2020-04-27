package stencil_test

import (
	"log"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/argots/stencil/pkg/stencil"
)

func TestCopyFile(t *testing.T) {
	//nolint: lll
	code := `{{ stencil.CopyFile "test1" "/test.txt" "source" }}`
	src := `{{ stencil.OS }}`
	var got []byte
	fs := fakeFS{
		files: map[string]string{"fake.stencil": code, "source": src},
		write: func(name string, data []byte, mode os.FileMode) error {
			if strings.HasSuffix(name, "/test.txt") {
				got = data
			}
			return nil
		},
	}

	discard := discardLogger{}
	if err := stencil.New(discard, discard, nil, fs).Run("fake.stencil"); err != nil {
		t.Fatal("Pull", err)
	}
	if string(got) != runtime.GOOS {
		t.Error("Got", string(got), "\nExpected", runtime.GOOS)
	}
}

type fakeFS struct {
	files map[string]string
	write func(fname string, data []byte, mode os.FileMode) error
}

func (f fakeFS) Write(fname string, data []byte, mode os.FileMode) error {
	if f.write == nil {
		return nil
	}
	return f.write(fname, data, mode)
}

func (f fakeFS) Read(fname string) ([]byte, error) {
	if v, ok := f.files[fname]; ok {
		return []byte(v), nil
	}

	discard := discardLogger{}
	fs := &stencil.FS{Verbose: discard, Errorl: discard}
	return fs.Read(fname)
}

func (f fakeFS) Remove(path string) error {
	return nil
}

func (f fakeFS) RemoveAll(path string) error {
	return nil
}

type discardLogger struct{}

func (discardLogger) Printf(fmt string, v ...interface{}) {
	log.Printf(fmt, v...)
}
