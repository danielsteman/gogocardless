package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	DBURL     string `json:"db_url"`
}

func LoadAppConfig(path string) (AppConfig, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return AppConfig{
		SecretID:  os.Getenv("SECRET_ID"),
		SecretKey: os.Getenv("SECRET_KEY"),
		DBURL:     os.Getenv("DB_URL"),
	}, nil
}
