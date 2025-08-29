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

	// select logger
	logger := gorm_logger.Default.LogMode(gorm_logger.Error)
	switch os.Getenv("DB_LOG_LEVEL") {
	case "info":
		logger = gorm_logger.Default.LogMode(gorm_logger.Info)
	case "silent":
		logger = gorm_logger.Default.LogMode(gorm_logger.Silent)
	}

	dbprefix := os.Getenv("DB_PREFIX")
	host := os.Getenv(dbprefix + "DB_HOST")
	user := os.Getenv(dbprefix + "DB_USERNAME")
	password := os.Getenv(dbprefix + "DB_PASSWORD")
	dbname := os.Getenv(dbprefix + "DB_NAME")
	port := os.Getenv(dbprefix + "DB_PORT")

	// connet to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger,
	})
	if err != nil {
		if host == "" {
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
