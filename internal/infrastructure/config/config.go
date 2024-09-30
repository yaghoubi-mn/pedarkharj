package config

import "time"

const (
	JWtRefreshExpireMinutes time.Duration = 30 * 24 * 60
	JWTAccessExpireMinutes  time.Duration = 15
	Debug                                 = true
)
