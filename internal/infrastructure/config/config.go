package config

import "time"

const (
	JWtRefreshExpire = 30 * 24 * time.Hour
	JWTAccessExpire  = 15 * time.Minute
	Debug            = true
)

var (
	BcryptCost = 16
)

func init() {
	if Debug {
		BcryptCost = 1
	}
}
