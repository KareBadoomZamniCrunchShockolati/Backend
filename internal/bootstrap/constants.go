package bootstrap

import (
	"time"
)

const (
	EnvFilePath    = "../../.env"
	DefaultAppPort = "8080"
)

var (
	JWTTokenExpiry = 24 * time.Hour
	JWTIssuer      = "challenge-app"
)

