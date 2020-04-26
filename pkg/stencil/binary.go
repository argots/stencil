package stencil

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/bmatcuk/doublestar"
)

const dialTimeout = time.Second * 10
const tlsTimeout = time.Second * 5
const httpTimeout = time.Second * 30
const targz = ".tar.gz"

// Binary implements managing binaries.
type Binary struct {
	*Stencil
}

// CopyFromArchive copies a file from an archive at the url.
// CopyFromArchive supports .tar, .tar.gz and .zip extensions for the archive.
func (b *Binary) CopyFromArchive(key, destination, url, file string) error {
	if b.Objects.existsArchiveFile(key, destination, url, file) {
		return nil
	}
	b.Objects.addArchiveFile(key, destination, url, file)
	seen := false
	err := b.extract(url, func(fname string, r func() io.ReadCloser) error {
		if !strings.EqualFold(file, fname) {
			return nil
		}

		src := r()
		defer src.Close()
		return b.copy(key+fname, destination, src)
	})
	if err == nil && !seen {
		err = errors.New("no such file: " + file)
	}
	return err
}

// CopyManyFromArchive extracts multiple files from an archive at the url.
// CopyManyFromArchive supports .tar, .tar.gz and .zip extensions for
// the archive. The glob pattern can be used to specify what files
// need to be extracted. See https://github.com/bmatcuk/doublestar for
// the set of allowed glob patterns. The destination is considered a
// folder.
func (b *Binary) CopyManyFromArchive(key, destination, url, glob string) error {
	if b.Objects.existsArchiveGlob(key, destination, url, glob) {
		return nil
	}
	b.Objects.addArchiveGlob(key, destination, url, glob)
	return b.extract(url, func(fname string, r func() io.ReadCloser) error {
		if match, err := doublestar.Match(glob, fname); err != nil || !match {
			return err
		}

		src := r()
		defer src.Close()
		return b.copy(key+fname, filepath.Join(destination, fname), src)
	})
}

func (b *Binary) extract(url string, visit func(name string, r func() io.ReadCloser) error) error {
	client := &http.Client{
		Timeout: httpTimeout,
		Transport: &http.Transport{
			Dial:                (&net.Dialer{Timeout: dialTimeout}).Dial,
			TLSHandshakeTimeout: tlsTimeout,
		},
	}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("http.Status " + resp.Status)
	}

	switch b.guessExtension(resp.Header.Get("Content-Type"), url) {
	case ".tar":
		return Untar(resp.Body, visit)
	case targz:
		r, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		return Untar(r, visit)
	case ".zip":
		return Unzip(resp.Body, visit)
	}

	return errors.New("Unknown destination URL extension " + url)
}

func (b *Binary) guessExtension(contentType, url string) string {
	switch contentType {
	case "application/zip":
		return ".zip"
	case "application/x-gzip":
		return targz
	}
	url = strings.ToLower(url)
	if strings.HasSuffix(url, targz) {
		return targz
	}
	return filepath.Ext(url)
}

func (b *Binary) copy(key, dest string, src io.Reader) error {
	_ = key
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}
	return b.Write(dest, data, 0766)
}
