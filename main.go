package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/yaghoubi-mn/pedarkharj/internal/user"
	"github.com/yaghoubi-mn/pedarkharj/pkg/cache"
	"github.com/yaghoubi-mn/pedarkharj/pkg/validator"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Println("WARNING: Cannot load env variables: ", err.Error())
	}

	// setup database
	db := SetupGrom()

	// setup cache
	cacheRepo := cache.New(db)

	// setup validator
	validatorIns := validator.NewValidator()

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
	log.Println("listening at :1111")
	log.Fatal(http.ListenAndServe(":1111", mux))
}

func SetupGrom() *gorm.DB {
	// connet to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran", os.Getenv("DB_HOST"), os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("ERROR: Cannot connect to database: ", err.Error())
	}

	err = db.AutoMigrate(
		&user.User{},
	)
	if err != nil {
		log.Println("WARNING: Cannot migrate tables: ", err.Error())
	}

	return db
}
