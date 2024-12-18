package config

import "time"

const (
	// debug mode
	Debug = true

	// s3
	AvatarPath = "avatars/"
)

var (
	// bcrypt cost
	BcryptCost = 16

	// jwt
	JWtRefreshExpire = 30 * 24 * time.Hour
	JWTAccessExpire  = 15 * time.Minute

	VerifyNumberCacheExpireTimeForNumberDelay = 3 * time.Minute
	VerifyNumberCacheExpireTime               = 10 * time.Minute
)

func init() {
	if Debug {
		BcryptCost = 1
		JWTAccessExpire = 24 * time.Hour
		VerifyNumberCacheExpireTimeForNumberDelay = 15 * time.Second
	}
}
