package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database `yaml:"db"`
	HTTP     `yaml:"http"`
}

type Database struct {
	DSN string `yaml:"dsn"`
}

type HTTP struct {
	PORT string `yaml:"port"`
}

func NewConfig() (*Config, error) {
	var filePath string
	flag.StringVar(&filePath, "c", "./configs/todos.yaml", "set config path")
	flag.Parse()

	c, err := parse(filePath)

	if err != nil {
		return nil, err
	}

	return c, nil
}

func parse(filepath string) (*Config, error) {
	f, err := os.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	var c Config

	if err := yaml.Unmarshal(f, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
