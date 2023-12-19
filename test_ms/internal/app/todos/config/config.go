package config

type Config struct {
	Database
	HTTP
}

func NewConfig() (*Config, error) {
	return &Config{}, nil
}

type Database struct {
	DSN string
}

type HTTP struct {
	PORT string
}
