package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/yaghoubi-mn/pedarkharj/internal/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/cache"
	"github.com/yaghoubi-mn/pedarkharj/pkg/jwt"
	"github.com/yaghoubi-mn/pedarkharj/pkg/validator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
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

	// create router
	muxV1 := http.NewServeMux()

	// user setup
	userRepo := user.NewGormUserRepository(db)
	userService := user.NewUserService(userRepo, cacheRepo, &validatorIns)
	userHandler := user.NewHandler(userService)
	user.Route("/users", muxV1, userHandler)

	// mux.Handle("/", middleware(fun))
	mux := http.NewServeMux()

	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", muxV1))

	slog.Info("listening at http://127.0.0.1:1111")
	slog.Error(http.ListenAndServe(":1111", mux).Error())
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
		&user.User{},
	)
	if err != nil {
		slog.Warn("Cannot migrate tables", "error", err.Error())
	}

	return db
}
