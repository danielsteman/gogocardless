package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDB(DBName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=localhost user=admin password=admin dbname=%s port=5432 sslmode=disable TimeZone=Europe/Amsterdam", DBName)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return db, nil
}
