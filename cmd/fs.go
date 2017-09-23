package main

import (
	"io/ioutil"
	"os"
)

type RealFS struct{}

func (f RealFS) Mkdir(path string, perms os.FileMode) error {
	return os.MkdirAll(path, perms)
}
func (f RealFS) WriteFile(path string, data []byte, perms os.FileMode) error {
	return ioutil.WriteFile(path, data, perms)
}
func (f RealFS) ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
