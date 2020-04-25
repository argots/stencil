package stencil

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// NewCache returns a new cache.
func NewCache(f *flag.FlagSet, stdin io.Reader, stdout io.Writer) (*Cache, error) {
	c := &Cache{Vars: NewVars(f, stdin, stdout)}
	return c, c.init()
}

// Cache keeps track of existing keys and associated files.
type Cache struct {
	BaseDir string
	Keys    []string
	*Vars
}

func (c *Cache) init() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	c.BaseDir = path
	for dir := path; dir != "/"; dir = filepath.Dir(dir) {
		s := filepath.Join(dir, ".stencil")
		if fi, err := os.Stat(s); err == nil && fi.IsDir() {
			c.BaseDir = dir
			return c.loadKeys()
		}
	}

	return c.loadKeys()
}

func (c *Cache) loadKeys() error {
	data, err := ioutil.ReadFile(filepath.Join(c.BaseDir, ".stencil", "keys.json"))
	if os.IsNotExist(err) {
		return nil
	}
	if err == nil {
		err = json.Unmarshal(data, &c.Keys)
	}
	return err
}

func (c *Cache) saveKeys() error {
	path := filepath.Join(c.BaseDir, ".stencil", "keys.json")
	if err := os.MkdirAll(filepath.Dir(path), 0766); err != nil {
		return err
	}
	data, err := json.Marshal(c.Keys)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0666)
}

// Get fetches a file from cache by key.
func (c *Cache) Get(key string) ([]byte, error) {
	path := filepath.Join(c.BaseDir, ".stencil", "cache", key, "latest.txt")
	return ioutil.ReadFile(path)
}

// Set sets a file from cache by key.
func (c *Cache) Set(key string, data []byte) error {
	path := filepath.Join(c.BaseDir, ".stencil", "cache", key, "latest.txt")
	if err := os.MkdirAll(filepath.Dir(path), 0766); err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, data, 0666); err != nil {
		return err
	}
	c.Keys = append(c.Keys, key)
	return c.saveKeys()
}

// Remove removes a file from cache by key.
func (c *Cache) Remove(key string) error {
	for idx, k := range c.Keys {
		if k != key {
			continue
		}
		path := filepath.Join(c.BaseDir, ".stencil", "cache", key, "latest.txt")
		if err := os.RemoveAll(path); err != nil {
			return err
		}
		c.Keys = append(c.Keys[:idx], c.Keys[idx+1:]...)
		return c.saveKeys()
	}
	return nil
}
