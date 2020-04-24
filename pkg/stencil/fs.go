package stencil

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

// MaxFileSize specifies the max file size in the cache.
const MaxFileSize = 1000000

// FS implements a FileSystem interface.
type FS struct {
	Verbose, Errorl Logger
}

// Write saves a file within the local directory
func (fs *FS) Write(path string, data []byte, mode os.FileMode) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	path = filepath.Join(cwd, filepath.Clean(path))
	return ioutil.WriteFile(path, data, mode)
}

// Read reads the contents of the path and returns them as bytes.
func (fs *FS) Read(path string) ([]byte, error) {
	u, err := url.Parse(path)
	if err == nil && u.Scheme == "git" {
		if u.Opaque != "" {
			u.Host = u.Opaque
		}
		return fs.ReadGit(u.Host, u.Path, u.Fragment)
	}
	if err == nil && u.Scheme != "" {
		return nil, errors.New("unsupport path format: " + path)
	}
	return ioutil.ReadFile(path)
}

func (fs *FS) ReadGit(host, path, fragment string) ([]byte, error) {
	branch := "master"
	gitURL := ""

	if fragment != "" {
		if idx := strings.Index(fragment, "/"); idx >= 0 {
			path += fragment[idx:]
			fragment = fragment[:idx]
		}
	}

	if fragment != "" {
		branch = fragment
	}
	gitURL = host + path
	path = gitURL[strings.Index(gitURL, ".git/")+4:]
	gitURL = gitURL[:strings.Index(gitURL, ".git/")+4]

	objectsDir, err := ioutil.TempDir("", "stencil-git-objects")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(objectsDir)
	workDir, err := ioutil.TempDir("", "stencil-git-work")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(workDir)

	objs := filesystem.NewStorage(osfs.New(objectsDir), cache.NewObjectLRU(MaxFileSize))
	work := osfs.New(workDir)
	_, err = git.Clone(objs, work, &git.CloneOptions{
		URL:           gitURL,
		ReferenceName: plumbing.ReferenceName("refs/heads/" + branch),
		SingleBranch:  true,
		Progress:      progress{fs},
		Tags:          git.AllTags,
	})
	if err != nil {
		fs.Errorl.Printf("git clone %s %v\n", gitURL, err)
		return nil, err
	}
	f, err := work.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

type progress struct {
	*FS
}

func (p progress) Write(data []byte) (int, error) {
	return len(data), nil
}
