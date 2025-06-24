package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/yaghoubi-mn/pedarkharj/docs"
	app_device "github.com/yaghoubi-mn/pedarkharj/internal/application/device"
	app_user "github.com/yaghoubi-mn/pedarkharj/internal/application/user"
	domain_device "github.com/yaghoubi-mn/pedarkharj/internal/domain/device"
	domain_user "github.com/yaghoubi-mn/pedarkharj/internal/domain/user"
	_ "github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/config"
	gorm_repository "github.com/yaghoubi-mn/pedarkharj/internal/infrastructure/repository/gorm"
	interfaces_rest_v1 "github.com/yaghoubi-mn/pedarkharj/internal/interfaces/rest/v1"
	"github.com/yaghoubi-mn/pedarkharj/pkg/cache"
	"github.com/yaghoubi-mn/pedarkharj/pkg/database"
	"github.com/yaghoubi-mn/pedarkharj/pkg/datatypes"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/s3"
	"github.com/yaghoubi-mn/pedarkharj/pkg/validator"
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

	// setup s3
	s3.Init()

	// setup database
	db, err := database.SetupGrom()
	if err != nil {
		slog.Error("database", "error", err)
		os.Exit(1)
	}

	// setup cache
	cacheRepo := cache.New(db)

	// migrate
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		slog.Info("migrating tables")

		err := database.MigrateTables(
			db,
			domain_user.User{},
			domain_device.Device{},
		)

		if err != nil {
			slog.Warn("Cannot migrate tables", "error", err.Error())
		}
		err = cache.MigrateTables(db)
		if err != nil {
			slog.Warn("Cannot migrate tables", "error", err.Error())
		}
		return
	}

	// setup validator
	validatorIns := validator.NewValidator()

	// setup jwt
	if jwtSecretKey := os.Getenv("JWT_SECRET_KEY"); jwtSecretKey == "" {
		slog.Error("JWT_SECRET_KEY not found in ENV")
		os.Exit(1)
	} else {
		jwt.Init(jwtSecretKey)
	}

	mux := setupRouter(db, validatorIns, cacheRepo)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	if len(os.Args) > 1 && os.Args[1] == "debug" {
		slog.Info("listening at http://127.0.0.1:2222")
		slog.Error(http.ListenAndServe(":2222", mux).Error())
	} else {
		slog.Info("Swagger: http://127.0.0.1:8000/swagger/index.html")
		slog.Info("listening at http://127.0.0.1:8000")
		slog.Error(http.ListenAndServe(":8000", mux).Error())

	}
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
