package stencil

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// FileObj tracks a single file copied locally.
type FileObj struct {
	Loc, URL string
}

// FileArchiveObj tracks an archive.
type FileArchiveObj struct {
	Many           bool
	Loc, URL, File string
}

// Objects tracks a collection of objects
type Objects struct {
	*Stencil     `json:"-"`
	Before       *Objects `json:"-"`
	Pulls        map[string]bool
	Files        map[string]*FileObj
	FileArchives map[string]*FileArchiveObj
	Bools        map[string]bool
	Strings      map[string]string
}

// LoadObjects loads all the objects from the .stencil directory.
func (o *Objects) LoadObjects() error {
	data, err := o.Read(".stencil/objects.json")
	if err == nil {
		return json.Unmarshal(data, o.Before)
	}
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// SaveObjects saves all the objects to the .stencil directory.
func (o *Objects) SaveObjects() error {
	data, err := json.MarshalIndent(o, "", "  ") //nolint: staticcheck
	if err != nil {
		return err
	}

	return o.Write(".stencil/objects.json", data, 0666)
}

func (o *Objects) addPull(url string) {
	o.Pulls[url] = true
}

func (o *Objects) addFile(key, dest, url string) {
	o.Files[key] = &FileObj{dest, url}
}

func (o *Objects) addArchiveFile(key, dest, url, file string) {
	o.FileArchives[key] = &FileArchiveObj{false, dest, url, file}
}

func (o *Objects) addArchiveGlob(key, dest, url, glob string) {
	o.FileArchives[key] = &FileArchiveObj{true, dest, url, glob}
}

func (o *Objects) existsArchiveFile(key, dest, url, file string) bool {
	if f, ok := o.FileArchives[key]; ok && !f.Many {
		return f.Loc == dest && f.URL == url && f.File == file
	}
	return false
}

func (o *Objects) existsArchiveGlob(key, dest, url, file string) bool {
	if f, ok := o.FileArchives[key]; ok && f.Many {
		return f.Loc == dest && f.URL == url && f.File == file
	}
	return false
}

// GC removes any files and dirs that are no longer active.
func (o *Objects) GC() error {
	files := map[string]bool{}
	dirs := map[string]bool{}
	o.Before.visitFile(func(file string) {
		files[file] = true
	})
	o.Before.visitDir(func(dir string) {
		dirs[dir] = true
	})
	o.visitFile(func(file string) {
		delete(files, file)
		o.deleteParents(dirs, file)
	})
	o.visitDir(func(dir string) {
		delete(dirs, dir)
		o.deleteParents(dirs, dir)
	})

	for file := range files {
		if err := o.Remove(file); err != nil {
			return err
		}
	}
	for dir := range dirs {
		if err := o.RemoveAll(dir); err != nil {
			return err
		}
	}
	return nil
}

func (o *Objects) deleteParents(dirs map[string]bool, file string) {
	for dir := filepath.Dir(file); dir != file; dir = filepath.Dir(file) {
		delete(dirs, dir)
		file = dir
	}
}

func (o *Objects) visitFile(fn func(file string)) {
	for _, f := range o.Files {
		fn(filepath.Clean(f.Loc))
	}
	for _, f := range o.FileArchives {
		fn(filepath.Clean(f.Loc))
	}
}

func (o *Objects) visitDir(fn func(dir string)) {
	for _, f := range o.FileArchives {
		if f.Many {
			fn(filepath.Clean(f.Loc))
		}
	}
}
