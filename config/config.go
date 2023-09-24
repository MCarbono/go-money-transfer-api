package config

import (
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	DBHost string
}

func LoadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.DBHost = os.Getenv("DB_HOST")
	return cfg, nil
}
