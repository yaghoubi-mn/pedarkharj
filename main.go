package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger/v2"
	_ "github.com/yaghoubi-mn/pedarkharj/docs"
	app_device "github.com/yaghoubi-mn/pedarkharj/internal/application/device"
	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	_ "github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/config"
	gorm_repository "github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/repository/gorm"
	interfaces_rest_v1 "github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1"
	"github.com/yaghoubi-mn/pedarkharj/pkg/cache"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/validator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title Pedarkharj
// @version 1.0.0
// @description Pedarkharj project
// @host localhost:1111
// @BasePath /api/v1
func main() {
	// setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	// load .env variables
	err := godotenv.Load()
	if err != nil {
		slog.Warn("Cannot load env variables", "error", err.Error())
	}

	// setup database
	db := SetupGrom()

	// setup cache
	cacheRepo := cache.New(db)

	// setup validator
	validatorIns := validator.NewValidator()

	// setup jwt
	if jwtSecretKey := os.Getenv("JWT_SECRET_KEY"); jwtSecretKey == "" {
		slog.Error("JWT_SECRET_KEY not found in ENV")
	} else {
		jwt.Init(jwtSecretKey)
	}

	mux := setupRouter(db, validatorIns, cacheRepo)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	slog.Info("listening at http://127.0.0.1:1111")
	slog.Error(http.ListenAndServe(":1111", mux).Error())
}

func setupRouter(db *gorm.DB, validatorIns datatypes.Validator, cacheRepo datatypes.CacheRepository) *http.ServeMux {

	// setup domain
	userDomainService := domain_user.NewUserService(validatorIns)
	deviceDomainService := domain_device.NewDeviceService(validatorIns)

	// setup repository
	userRepo := gorm_repository.NewGormUserRepository(db)
	deviceRepo := gorm_repository.NewGormDeviceRepository(db)

	// setup application
	deviceAppService := app_device.NewDeviceAppService(deviceRepo, deviceDomainService)
	userAppService := app_user.NewUserService(userRepo, cacheRepo, deviceAppService, userDomainService)

	// setup router
	muxV1 := interfaces_rest_v1.NewRouter(userAppService, deviceAppService)

	return muxV1
}

func SetupGrom() *gorm.DB {
	// connet to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Cannot connect to database", "error", err.Error())
		os.Exit(1)
	}

	err = db.AutoMigrate(
		&domain_user.User{},
		&domain_device.Device{},
	)
	if err != nil {
		slog.Warn("Cannot migrate tables", "error", err.Error())
	}

	return db
}
