package config

import (
	"log/slog"
	"os"
	"time"
)

const (

	// s3
	AvatarPath = "avatars/"
)

var (
	Debug = true

	// bcrypt cost
	BcryptCost = 16

	// jwt
	JWtRefreshExpire = 30 * 24 * time.Hour
	JWTAccessExpire  = 15 * time.Minute

	VerifyNumberCacheExpireTimeForNumberDelay = 3 * time.Minute
	VerifyNumberCacheExpireTime               = 10 * time.Minute
)

func init() {
	if os.Getenv("DEBUG") == "false" {
		Debug = false
		slog.Info("---Debug mode is off---")
	}

	if Debug {
		BcryptCost = 1
		JWTAccessExpire = 24 * time.Hour
		VerifyNumberCacheExpireTimeForNumberDelay = 15 * time.Second
	}
}
