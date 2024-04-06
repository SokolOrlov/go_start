package config

import (
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

func Read() (*Config, error) {
	var filePath string
	flag.StringVar(&filePath, "config", "./cmd/client/config/config.yaml", "set config path")
	flag.Parse()

	// if filePath == "" {
	// 	filePath = os.Getenv("CONFIG_PATH")
	// }

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
