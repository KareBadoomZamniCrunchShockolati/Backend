package bootstrap

import (
	"time"
)

const (
	EnvFilePath = "../../.env"
	AppPort     = "8080"
	RedisPort   = "6379"
)

var (
	JWTTokenExpiry = 24 * time.Hour
	JWTIssuer      = "challenge-app"
)
