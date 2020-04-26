package stencil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/argots/stencil/pkg/stencil"
)

func TestBaseDirNoRoot(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd", err)
	}

	baseDir, err := stencil.BaseDir()
	if err != nil {
		t.Fatal("stencil.BaseDir", err)
	}

	if baseDir != dir {
		t.Error("Unexpected base dir", baseDir, dir)
	}
}

func TestBaseDirRoot(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd", err)
	}

	defer os.Chdir(dir) //nolint: errcheck

	somedir := filepath.Join(dir, "testdata/root/somedir")
	if err = os.Chdir(somedir); err != nil {
		t.Fatal("os.Chdir", err)
	}

	baseDir, err := stencil.BaseDir()
	if err != nil {
		t.Fatal("stencil.BaseDir", err)
	}

	root := filepath.Join(dir, "testdata/root")
	if baseDir != root {
		t.Error("Unexpected base dir", baseDir, root)
	}
}
