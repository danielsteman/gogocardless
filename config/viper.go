package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	DbURL        string `mapstructure:"DB_URL"`
	DbDriver     string `mapstructure:"DB_DRIVER"`
	ServeAddress string `mapstructure:"SERVE_ADDRESS"`
}

func LoadAppConfig() (AppConfig, error) {
	path := "settings.yml"
	viper.AddConfigPath(path)
	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		return AppConfig{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return config, nil
}
