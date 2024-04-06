package config

type Config struct {
	DB    DB    `yaml:"db"`
	KAFKA KAFKA `yaml:"kafka"`
}

type DB struct {
	USER          string `yaml:"user"`
	PWD           string `yaml:"pwd"`
	MigrationPath string `yaml:"migrationPath"`
}

type KAFKA struct {
	BROKER   BROKER   `yaml:"broker"`
	CONSUMER CONSUMER `yaml:"consumer"`
	PRODUCER PRODUCER `yaml:"producer"`
}
type BROKER struct {
	URL  string `yaml:"url"`
	PORT string `yaml:"port"`
}

type CONSUMER struct {
	TOPIC    string `yaml:"topic"`
	GROUP    string `yaml:"group"`
	ASSIGNOR string `yaml:"assignor"`
	OLDEST   bool   `yaml:"oldest"`
}
type PRODUCER struct {
	CLIENTTOPIC string `yaml:"clientTopic"`
	DRIVERTOPIC string `yaml:"driverTopic"`
}
