package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	SecretID  string `mapstructure:"SECRET_ID"`
	SecretKey string `mapstructure:"SECRET_KEY"`
}

func LoadAppConfig(path string) (AppConfig, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.SetConfigFile(path)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		return AppConfig{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}
