package config

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type ToBackup = string
type Manifest = string

type Config struct {
	Local struct {
		Paths map[ToBackup]Manifest `yaml:"paths"`
	} `yaml:"local"`
}

func (c *Config) validate() error {
	for _, path := range c.Local.Paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func NewConfig(path string) (*Config, error) {
	fd, err := os.Open(path)
	if err != nil {
		fmt.Println("Error while opening config file")
		return nil, err
	}
	defer fd.Close()

	var cfg Config
	decoder := yaml.NewDecoder(fd)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.validate()

	return &cfg, err
}
