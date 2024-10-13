package config

import "time"

const (
	Debug = true

	// jwt
	JWtRefreshExpire = 30 * 24 * time.Hour
	JWTAccessExpire  = 15 * time.Minute

	// s3
	AvatarPath = "root/avatars/"
)

var (
	BcryptCost = 16
)

func init() {
	if Debug {
		BcryptCost = 1
	}
}
