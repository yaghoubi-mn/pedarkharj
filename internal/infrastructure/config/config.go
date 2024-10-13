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
)

func init() {
	if Debug {
		BcryptCost = 1
		JWTAccessExpire = 24 * time.Hour
	}
}
