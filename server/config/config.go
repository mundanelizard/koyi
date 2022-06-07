package config

import "time"

/**
Configure as environment variable?

UI - domain/ip address
*/

const (
	RefreshTokenCookieMaxAge = 200
	ServerDomain             = "localhost:5002"
	IsProduction             = false
	AccessTokenDuration      = time.Hour * time.Duration(10) // 10 hours
	RefreshTokenDuration     = time.Hour * time.Duration(20) // 20 hours
	AccessTokenSecretKey     = "super-secret-key"
	RefreshTokenSecretKey    = "super-duper-secret-key"
	MongoUri                 = "mongodb://localhost:27017/"
	AverageServerTimeout     = 20 * time.Second
	UserDatabaseName         = "users"
	BackgroundTaskTimeout    = 60 * time.Second
	JWTIssuerName            = "koyi"
	JWTAudience              = "koyi-client"
	EmailAddress             = "info@koyi.com"
	ValidateNewDevice        = true
)
