package cmd

import (
	"io/ioutil"
	"os"
)

// RealFS is an implementation of cache.FS
type RealFS struct{}

// Mkdir implementation
func (f RealFS) Mkdir(path string, perms os.FileMode) error {
	return os.MkdirAll(path, perms)
}

// WriteFile implementation
func (f RealFS) WriteFile(path string, data []byte, perms os.FileMode) error {
	return ioutil.WriteFile(path, data, perms)
}

// ReadFile implementation
func (f RealFS) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
