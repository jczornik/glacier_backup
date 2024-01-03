package config

import (
	"fmt"
	"os"

	"github.com/go-yaml/yaml"
)

type BackupSrc = string
type BackupDst = string

type Config struct {
	Local struct {
		Paths map[BackupSrc]BackupDst `yaml:"paths"`
	} `yaml:"local"`
}

func (c *Config) validate() error {
	for src := range c.Local.Paths {
		if _, err := os.Stat(src); os.IsNotExist(err) {
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
