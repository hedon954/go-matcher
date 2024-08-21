package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// FileLoader loads the config from file.
type FileLoader struct {
	path string

	sync.RWMutex
	config *Config
}

func NewFileLoader(path string) *FileLoader {
	fl := &FileLoader{path: path}
	fl.load()
	return fl
}

func (fl *FileLoader) Get() *Config {
	fl.RLock()
	defer fl.RUnlock()
	return fl.config
}

func (fl *FileLoader) load() {
	config, err := load(fl.path)
	if err != nil {
		panic(err)
	}

	fl.Lock()
	defer fl.Unlock()
	fl.config = config
}

// load loads the config from file.
func load(path string) (*Config, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %w", err)
	}
	c := &Config{}
	err = yaml.Unmarshal(bs, c)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config error: %w", err)
	}
	return c, nil
}
