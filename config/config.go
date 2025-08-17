package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Debug            bool     `default:"true" split_words:"true"`
	Port             int      `default:"8080" split_words:"true"`
	PsqlDB           Psql     `split_words:"true"`
	MongoDB          Mongo    `split_words:"true"`
	AcceptedVersions []string `required:"true" split_words:"true"`
}

// Database represents the configuration for the database connection.
type Psql struct {
	Driver    string
	User      string
	Password  string
	Port      int
	Host      string
	DATABASE  string
	Schema    string
	MaxActive int
	MaxIdle   int
}

type Mongo struct {
	URI      string
	User     string
	Password string
	Database string
	MaxPool  uint64
}

// LoadConfig loads the configuration for the application based on the given appName.
func Load(appName string) (*EnvConfig, error) {
	var cfg EnvConfig

	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	err := envconfig.Process(appName, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
