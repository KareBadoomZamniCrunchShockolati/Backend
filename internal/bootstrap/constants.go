package bootstrap

import "time"

const (
	EnvFilePath = "../../.env"
	AppPort     = "8081"
	RedisPort   = "6379"
)

var (
	JWTTokenExpiry = 24 * time.Hour
	JWTIssuer      = "challenge-app"
)

const (
	ProfileMinWidth  = 320
	ProfileMinHeight = 320
	PostMinDimension = 1080
	MaxProfileSize   = 5 << 20  // 5 MB
	MaxPostImageSize = 10 << 20 // 10 MB
	MaxCoverSize     = 10 << 20 // 10 MB 
	CoverMinWidth  = 1080
	CoverMinHeight = 540
)

const (
	PostImagesTTL = 15 * time.Minute
	MaxPostImages = 10
)

