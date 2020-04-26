package stencil

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/semver"

	"github.com/go-git/go-billy/v5"
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
	BaseDir         string
	Verbose, Errorl Logger
}

// Remove removes a file.
func (fs *FS) Remove(path string) error {
	fs.Verbose.Printf("Deleting file %s\n", path)
	return os.Remove(path)
}

// RemoveAll removes a directory and all its contents.
func (fs *FS) RemoveAll(path string) error {
	fs.Verbose.Printf("Deleting dir %s\n", path)
	return os.RemoveAll(path)
}

// Write saves a file within the local directory
func (fs *FS) Write(path string, data []byte, mode os.FileMode) error {
	path = filepath.Join(fs.BaseDir, filepath.Clean(path))
	if err := os.MkdirAll(filepath.Dir(path), 0766); err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, mode)
}

// Read reads the contents of the path and returns them as bytes.
func (fs *FS) Read(path string) ([]byte, error) {
	if _, _, gitPath, ok := fs.parseGitURL(path); ok {
		var result []byte
		err := fs.withGit(path, func(r *git.Repository, bfs billy.Filesystem) error {
			f, err := bfs.Open(gitPath)
			if err != nil {
				return err
			}
			defer f.Close()

			result, err = ioutil.ReadAll(f)
			return err
		})
		return result, err
	}
	return ioutil.ReadFile(filepath.Join(fs.BaseDir, filepath.Clean(path)))
}

// Resolve pins a git path to a specific hash/commit.
func (fs *FS) Resolve(path string) (string, error) {
	branch, cloneURL, path, ok := fs.parseGitURL(path)
	if !ok {
		return "", errors.New("not a git url: " + path)
	}

	err := fs.withGit(path, func(r *git.Repository, _ billy.Filesystem) error {
		head, err := r.Head()
		if err != nil {
			return err
		}

		branch = head.Name().Short()
		return nil
	})
	if err != nil {
		return "", err
	}
	return "git:" + cloneURL + "#" + branch + path, nil
}

func (fs *FS) parseGitURL(gitURL string) (branch, cloneURL, path string, ok bool) {
	u, err := url.Parse(gitURL)
	if err != nil || u.Scheme != "git" {
		return "", "", "", false
	}
	if u.Opaque != "" {
		u.Host = u.Opaque
	}

	branch = "master"
	cloneURL = ""
	path = u.Path

	if u.Fragment != "" {
		if idx := strings.Index(u.Fragment, "/"); idx >= 0 {
			path += u.Fragment[idx:]
			u.Fragment = u.Fragment[:idx]
		}
	}

	if u.Fragment != "" {
		branch = u.Fragment
	}
	cloneURL = u.Host + path
	path = cloneURL[strings.Index(cloneURL, ".git/")+4:]
	cloneURL = cloneURL[:strings.Index(cloneURL, ".git/")+4]
	return branch, cloneURL, path, true
}

// withGit calls te provided function on the git directory.
func (fs *FS) withGit(url string, fn func(r *git.Repository, fs billy.Filesystem) error) error {
	branch, cloneURL, _, _ := fs.parseGitURL(url)

	objectsDir, err := ioutil.TempDir("", "stencil-git-objects")
	if err != nil {
		return err
	}
	defer os.RemoveAll(objectsDir)
	workDir, err := ioutil.TempDir("", "stencil-git-work")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workDir)

	objs := filesystem.NewStorage(osfs.New(objectsDir), cache.NewObjectLRU(MaxFileSize))
	work := osfs.New(workDir)
	r, err := git.Clone(objs, work, &git.CloneOptions{
		URL:           cloneURL,
		ReferenceName: fs.refName(branch),
		SingleBranch:  true,
		Progress:      progress{fs},
		Tags:          git.AllTags,
	})
	if err != nil {
		fs.Errorl.Printf("git clone %s %v\n", cloneURL, err)
		return err
	}
	return fn(r, work)
}

func (fs *FS) refName(branch string) plumbing.ReferenceName {
	if semver.IsValid(branch) {
		return plumbing.NewTagReferenceName(branch)
	}
	return plumbing.NewBranchReferenceName(branch)
}

type progress struct {
	*FS
}

func (p progress) Write(data []byte) (int, error) {
	return len(data), nil
}
