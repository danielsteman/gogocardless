package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	SecretID    string `json:"secret_id"`
	SecretKey   string `json:"secret_key"`
	DBURL       string `json:"db_url"`
	DBName      string `json:"db_name"`
	Port        int    `json:"port"`
	RedirectURL string `json:"redirect_url"`
	JWTSecret   string `json:"jwt_secret"`
}

var (
	Config AppConfig
	once   sync.Once
)

func LoadAppConfig(path string) {
	once.Do(func() {
		err := godotenv.Load(path)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error loading config from %s", path))
		}
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			log.Fatal(fmt.Sprintf("PORT should be a number"))
		}
		Config = AppConfig{
			SecretID:    os.Getenv("SECRET_ID"),
			SecretKey:   os.Getenv("SECRET_KEY"),
			DBURL:       os.Getenv("DB_URL"),
			DBName:      os.Getenv("DB_NAME"),
			Port:        port,
			RedirectURL: os.Getenv("REDIRECT_URL"),
			JWTSecret:   os.Getenv("JWT_SECRET"),
		}
	})
}

func ResetAppConfig() {
	once = sync.Once{}
}
