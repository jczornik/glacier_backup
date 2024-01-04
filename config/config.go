package config

import (
	"errors"
	"os"
	"log"

	"github.com/go-yaml/yaml"
)

type BackupSrc = string
type BackupDst = string
type VaultName = string
type AwsSharedProfile = string

type BackupConfig struct {
	Src   BackupSrc `yaml:"src"`
	Dst   BackupDst `yaml:"dst"`
	Keep  bool      `yaml:"keep"`
	Vault VaultName `yaml:"vault"`
}

type AWSConfig struct {
	Profile   string `yaml:"profile"`
	AccountID string `yaml:"account"`
}

type Config struct {
	Backups []BackupConfig `yaml:"backup"`
	AWS     AWSConfig      `yaml:"aws"`
}

func (c *Config) validate() error {
	for _, backup := range c.Backups {
		if err := backup.validate(); err != nil {
			return err
		}
	}

	if err := c.AWS.validate(); err != nil {
		return err
	}

	return nil
}

func (c *BackupConfig) validate() error {
	if len(c.Src) == 0 || len(c.Dst) == 0 {
		return errors.New("Backup src and dst cannot be empty")
	}

	if _, err := os.Stat(c.Src); os.IsNotExist(err) {
		return err
	}

	if _, err := os.Stat(c.Dst); os.IsNotExist(err) {
		return err
	}

	if len(c.Vault) == 0 {
		return errors.New("Vault name cannot be empty")
	}

	return nil
}

func (c *AWSConfig) validate() error {
	if len(c.Profile) == 0 {
		return errors.New("AWS profile cannot be empty")
	}

	if len(c.AccountID) == 0 {
		return errors.New("AWS profile cannot be empty")
	}

	return nil
}

func NewConfig(path string) (*Config, error) {
	fd, err := os.Open(path)
	if err != nil {
		log.Println("Cannot open config file")
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
