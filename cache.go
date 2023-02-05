package iokit

import (
	"github.com/pkgs-go/fu"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
)

const cacheGoData = ".cache/pkgs-go/iokit-cache/"

var FullCacheDir string

func init() {
	homedir, _ := os.LookupEnv("HOME")
	usr, err := user.Current()
	if err == nil {
		homedir = usr.HomeDir
	}
	if homedir == "" {
		homedir = "/tmp"
	}
	FullCacheDir, _ = filepath.Abs(filepath.Join(homedir, cacheGoData))
}

func CacheDir(d string) string {
	r := fu.Ife(filepath.IsAbs(d), d, path.Join(FullCacheDir, d))
	_ = os.MkdirAll(r, 0777)
	return r
}

func CacheFile(f string) string {
	r := fu.Ife(filepath.IsAbs(f), f, path.Join(FullCacheDir, f))
	_ = os.MkdirAll(path.Dir(r), 0777)
	return r
}

type Cache string

func (c Cache) File() InputOutput {
	return File(CacheFile(string(c)))
}

func (c Cache) String() string {
	return CacheFile(string(c))
}

func (c Cache) Remove() (err error) {
	s := CacheFile(string(c))
	_, err = os.Stat(s)
	if err == nil {
		return os.Remove(s)
	}
	return nil
}

func (c Cache) RemoveNoError() {
	_ = c.Remove() // ignore error
}

func (c Cache) Defined() bool {
	return string(c) != ""
}

func (c Cache) Exists() bool {
	if c.Defined() {
		if st, err := os.Stat(c.Path()); err == nil && st.Mode().IsRegular() {
			return true
		}
	}
	return false
}

func (c Cache) Path() string {
	return CacheFile(string(c))
}

func (c Cache) Open() (io.ReadCloser, error) {
	return File(c.Path()).Open()
}

func (c Cache) MustOpen() StrictReader {
	rd, err := c.Open()
	if err != nil {
		panic(err)
	}
	return StrictReader{rd}
}
