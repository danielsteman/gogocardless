package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	DBURL     string `json:"db_url"`
	DBName    string `json:"db_name"`
	Port      int    `json:"port"`
}

func LoadAppConfig(path string) (AppConfig, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return AppConfig{}, err
	}
	return AppConfig{
		SecretID:  os.Getenv("SECRET_ID"),
		SecretKey: os.Getenv("SECRET_KEY"),
		DBURL:     os.Getenv("DB_URL"),
		DBName:    os.Getenv("DB_NAME"),
		Port:      port,
	}, nil
}
