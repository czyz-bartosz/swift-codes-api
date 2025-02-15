package configs

import (
	"awesomeProject/dbs"
	"os"
)

type Config struct {
	DBConfig dbs.Config
}

func GetConfig() *Config {
	config := Config{
		DBConfig: dbs.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Database: os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}
	return &config
}
