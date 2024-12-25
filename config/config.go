package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type SQLite3Config struct {
	Path   string `yaml:"path"`
	Schema string `yaml:"schema"`
}

type Config struct {
	API struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		LogLevel string `yaml:"log_level"`
		URL      string `yaml:"url"`
	} `yaml:"api"`
	Postgres PostgresConfig `yaml:"postgres"`
	SQLite3  SQLite3Config  `yaml:"sqlite3"`
}

func NewConfig(rootPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(rootPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
