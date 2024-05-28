package db

import (
	"fmt"

	"github.com/danielsteman/gogocardless/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	DBName string
	Port   int16
}

func GetDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=admin password=admin dbname=%s port=%d sslmode=disable TimeZone=Europe/Amsterdam", config.Config.DBName, config.Config.Port)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
