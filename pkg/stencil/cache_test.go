package stencil_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/argots/stencil/pkg/stencil"
)

func TestCacheNoRoot(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd", err)
	}

	c, err := stencil.NewCache()
	if err != nil {
		t.Fatal("stencil.NewCache", err)
	}

	if c.BaseDir != dir {
		t.Error("Unexpected base dir", c.BaseDir, dir)
	}
}

func TestCacheRoot(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd", err)
	}

	defer os.Chdir(dir) //nolint: errcheck

	somedir := filepath.Join(dir, "testdata/root/somedir")
	if err = os.Chdir(somedir); err != nil {
		t.Fatal("os.Chdir", err)
	}

	c, err := stencil.NewCache()
	if err != nil {
		t.Fatal("stencil.NewCache", err)
	}

	root := filepath.Join(dir, "testdata/root")
	if c.BaseDir != root {
		t.Error("Unexpected base dir", c.BaseDir, root)
	}
}

func TestCacheSetGet(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("os.Getwd", err)
	}

	defer os.Chdir(dir) //nolint: errcheck

	somedir := filepath.Join(dir, "testdata/root/somedir")
	if err = os.Chdir(somedir); err != nil {
		t.Fatal("os.Chdir", err)
	}

	c, err := stencil.NewCache()
	if err != nil {
		t.Fatal("stencil.NewCache", err)
	}

	defer os.RemoveAll(filepath.Join(dir, "testdata/root/.stencil/cache"))

	if err = c.Set("booya/22", []byte("hello world")); err != nil {
		t.Fatal("cache.Set", err)
	}

	if v, err := c.Get("booya/22"); err != nil || string(v) != "hello world" {
		t.Fatal("cache.Get", err, string(v))
	}

	if err = c.Remove("booya/22"); err != nil {
		t.Fatal("cache.Remove", err)
	}

	if _, err = c.Get("booya/22"); err == nil {
		t.Fatal("cache.Get")
	}
}
