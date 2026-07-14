package db

import (
	"fmt"
	"log"
	"time"

	"kiramidru/go-shopping/internal/auth"
	"kiramidru/go-shopping/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	gormLogLevel := logger.Warn
	if cfg.Env == "development" {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(gormLogLevel),
		TranslateError: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	psgDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	psgDB.SetMaxIdleConns(10)
	psgDB.SetMaxOpenConns(100)
	psgDB.SetConnMaxLifetime(1 * time.Hour)

	err = db.AutoMigrate(
		&auth.User{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	return db, nil
}
