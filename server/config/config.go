package config

import "time"

/**
Configure as environment variable?

UI - domain/ip address
*/

const (
	RefreshTokenCookieMaxAge = 200
	ServerDomain             = "google.com"
	IsProduction             = false
	AccessTokenDuration      = time.Hour * time.Duration(10) // 10 hours
	RefreshTokenDuration     = time.Hour * time.Duration(20) // 20 hours
	AccessTokenSecretKey     = "super-secret-key"
	RefreshTokenSecretKey    = "super-duper-secret-key"
)
