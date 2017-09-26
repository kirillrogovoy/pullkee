// Package cache provides an abstraction for data caching
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

var cachePath = path.Join(os.TempDir(), "pullk_cache", "cache.json")

// Cache is an interface for a Get/Set caching struct
type Cache interface {
	Set(string, interface{}) error
	Get(string, interface{}) (bool, error)
}

// FS is an interface for interacting with the file system
type FS interface {
	Mkdir(path string, perms os.FileMode) error
	WriteFile(path string, data []byte, perms os.FileMode) error
	ReadFile(path string) ([]byte, error)
}

// FSCache is an implementation of file-system cache using the FS interface
type FSCache struct {
	CachePath string
	FS        FS
}

// Set converts `x` to JSON and writes it to a file using `key`
func (c FSCache) Set(key string, x interface{}) error {
	err := c.FS.Mkdir(c.CachePath, 0744)
	if err != nil {
		return err
	}

	jsoned, err := json.Marshal(x)
	if err != nil {
		return err
	}

	return c.FS.WriteFile(c.filePath(key), jsoned, 0644)
}

// Get gets the contents of the file (if exists) using `key` and Unmarshals it to the struct `x`
func (c FSCache) Get(key string, x interface{}) (bool, error) {
	data, err := c.FS.ReadFile(c.filePath(key))
	if err != nil {
		// "no such file or directory" just should mean there's no cache entry
		if strings.Contains(err.Error(), "no such file or directory") {
			return false, nil
		}
		return false, err
	}

	return true, json.Unmarshal(data, x)
}

func (c FSCache) filePath(key string) string {
	return path.Join(c.CachePath, fmt.Sprintf("%s.json", key))
}
