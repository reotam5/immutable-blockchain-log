package internal

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dbInstance *gorm.DB

func InitDB() (*gorm.DB, error) {
	if dbInstance != nil {
		return dbInstance, nil
	}
	dsn := "host=localhost user=testuser password=testpass dbname=testdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.AutoMigrate(&LogEntry{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	dbInstance = db
	return dbInstance, nil
}
