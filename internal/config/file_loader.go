package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// FileLoader loads the config from file.
type FileLoader struct {
	path string
}

func NewFileLoader(path string) *FileLoader {
	return &FileLoader{path: path}
}

func (fl *FileLoader) Load() (*Config, error) {
	return load(fl.path)
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
