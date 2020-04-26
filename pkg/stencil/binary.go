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
)

const dialTimeout = time.Second * 10
const tlsTimeout = time.Second * 5
const httpTimeout = time.Second * 30
const targz = ".tar.gz"

// Binary implements managing binaries.
type Binary struct {
	FileSystem
}

// CopyFromArchive copies te binary at the provided URL.
// CopyFromArchive expects .tar, .tar.gz and .zip extensions.
func (b *Binary) CopyFromArchive(key, destination, url, file string) error {
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

	switch b.guessExtension(resp.Header.Get("Content-Type"), url) {
	case ".tar":
		return b.copy(key, destination, Untar(resp.Body, file))
	case targz:
		r, err := gzip.NewReader(resp.Body)
		if err != nil {
			return err
		}
		return b.copy(key, destination, Untar(r, file))
	case ".zip":
		r := Unzip(resp.Body, file)
		defer r.Close()
		return b.copy(key, destination, r)
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
