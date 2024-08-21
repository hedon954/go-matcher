package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// FileLoader loads the config from file.
type FileLoader[T any] struct {
	path string

	sync.RWMutex
	c *T
}

func NewFileLoader[T any](path string) *FileLoader[T] {
	fl := &FileLoader[T]{path: path}
	fl.load()
	return fl
}

func (fl *FileLoader[T]) Get() *T {
	fl.RLock()
	defer fl.RUnlock()
	return fl.c
}

func (fl *FileLoader[T]) load() {
	config, err := load[T](fl.path)
	if err != nil {
		panic(err)
	}

	fl.Lock()
	defer fl.Unlock()
	fl.c = config
}

// load loads the config from file.
func load[T any](path string) (*T, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}
	var c T
	err = yaml.Unmarshal(bs, &c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}
	return &c, nil
}
