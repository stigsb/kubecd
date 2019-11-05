package cache

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

func NewFileCache(scope string) *FileCache {
	dir := os.Getenv("KUBECD_CACHE")
	if dir == "" {
		me, _ := user.Current()
		dir = me.HomeDir
	}
	return &FileCache{directory: filepath.Join(dir, ".kubecd", "cache", "inspect")}
}

func NewMemoryCache(scope string) *MemoryCache {
	return &MemoryCache{cache: make(map[string][]byte)}
}

func Key(elements ...string) string {
	h := sha1.New()
	for _, el := range elements {
		h.Write([]byte(el))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

type Cache interface {
	Has(string) bool
	Get(string) []byte
	Set(string, []byte)
}

type FileCache struct {
	directory string
}

func (c FileCache) Has(key string) bool {
	path := filepath.Join(c.directory, key)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func (c FileCache) Get(key string) []byte {
	path := filepath.Join(c.directory, key)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

func (c FileCache) Set(key string, value []byte) {
	path := filepath.Join(c.directory, key)
	_ = ioutil.WriteFile(path, value, 0600)
}

var _ Cache = &FileCache{}

type MemoryCache struct {
	cache map[string][]byte
}

func (c MemoryCache) Has(key string) bool {
	_, has := c.cache[key]
	return has
}

func (c MemoryCache) Get(key string) []byte {
	return c.cache[key]
}

func (c MemoryCache) Set(key string, data []byte) {
	c.cache[key] = data
}

var _ Cache = &MemoryCache{}
