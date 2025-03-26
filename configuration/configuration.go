package configuration

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"os"
)

type Configuration struct {
	LogLevel         string `env:"LOG_LEVEL"`
	DatabaseUsername string `env:"DATABASE_USERNAME"`
	DatabasePassword string `env:"DATABASE_PASSWORD"`
	DatabaseHost     string `env:"DATABASE_HOST"`
	DatabasePort     string `env:"DATABASE_PORT"`
	DatabaseName     string `env:"DATABASE_NAME"`

	RedisHost            string `env:"REDIS_HOST"`
	RedisPort            string `env:"REDIS_PORT"`
	RedisPassword        string `env:"REDIS_PASSWORD"`
	RedisDatabase        int    `env:"REDIS_DATABASE"`
	MeiliSearchHost      string `env:"MEILISEARCH_HOST"`
	MeiliSearchPort      int    `env:"MEILISEARCH_PORT"`
	MeiliSearchMasterKey string `env:"MEILISEARCH_MASTER_KEY"`
}

func GetEnvConfig() (*Configuration, error) {
	var err error
	environment := os.Getenv("ENVIRONMENT")
	if environment != "" {
		err = godotenv.Load(".env." + environment)
	} else {
		err = godotenv.Load()
	}
	if err != nil {
		return nil, err
	}

	cfg := Configuration{
		LogLevel: "info",
	}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	return &cfg, nil
}
