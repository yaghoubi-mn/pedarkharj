package database

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

func SetupGrom() (*gorm.DB, error) {
	// connet to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gorm_logger.Default.LogMode(gorm_logger.Info),
	})
	if err != nil {
		if os.Getenv("DB_HOST") == "" {
			slog.Error("maybe env variable not loaded")
		}
		return nil, errors.New("Cannot connect to database: " + err.Error())
	}

	return db, nil
}

func MigrateTables(db *gorm.DB, models ...any) error {

	return db.AutoMigrate(models...)
	// &domain_user.User{},
	// &domain_device.Device{},
}

func SetupGromForTest() (*gorm.DB, error) {
	// connet to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_USERNAME"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gorm_logger.Default.LogMode(gorm_logger.Error),
	})
	if err != nil {
		if os.Getenv("TEST_DB_HOST") == "" {
			slog.Error("maybe env variable not loaded")
		}
		return nil, errors.New("Cannot connect to database: " + err.Error())
	}

	return db, nil
}
