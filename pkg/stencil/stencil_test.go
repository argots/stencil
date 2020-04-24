package stencil_test

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/argots/stencil/pkg/stencil"
)

func TestCopyFile(t *testing.T) {
	//nolint: lll
	code := `{{ stencil.CopyFile "test1" "/test.txt" "git:git@github.com:argots/stencil.git/pkg/stencil/testdata/test.stencil" }}`
	var got []byte
	fs := fakeFS{
		files: map[string]string{"fake.stencil": code},
		write: func(name string, data []byte, mode os.FileMode) error {
			if strings.HasSuffix(name, "/test.txt") {
				got = data
			}
			return nil
		},
	}

	expected, err := ioutil.ReadFile("testdata/test.stencil")
	if err != nil {
		t.Fatal(err)
	}

	discard := discardLogger{}
	if err = stencil.New(discard, discard, fs).Run("fake.stencil"); err != nil {
		t.Fatal("Pull", err)
	}
	if string(got) != string(expected) {
		t.Error("Got", string(got), "\nExpected", string(expected))
	}
}

type fakeFS struct {
	files map[string]string
	write func(fname string, data []byte, mode os.FileMode) error
}

func (f fakeFS) Write(fname string, data []byte, mode os.FileMode) error {
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

type discardLogger struct{}

func (discardLogger) Printf(fmt string, v ...interface{}) {
	log.Printf(fmt, v...)
}
