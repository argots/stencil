package stencil

import (
	"os"
	"path/filepath"
)

func BaseDir() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for dir := path; dir != "/"; dir = filepath.Dir(dir) {
		s := filepath.Join(dir, ".stencil")
		if fi, err := os.Stat(s); err == nil && fi.IsDir() {
			return dir, nil
		}
	}

	return path, nil
}
