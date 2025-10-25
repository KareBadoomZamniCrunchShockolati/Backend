package security

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"challenge-app/internal/bootstrap"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateToken(userID uint) (string, error)
	ValidateToken(signedToken string) (*Claims, error) 
}

type JwtServiceImpl struct {
	secretKey       []byte
	tokenExpiration time.Duration
	issuer          string
}

// NewJWTService creates a new JWTService implementation.
func NewJWTService(cfg *bootstrap.Env) *JwtServiceImpl {
	return &JwtServiceImpl{
		secretKey:       []byte(cfg.Security.JWTSecret),
		tokenExpiration: cfg.Security.TokenTTL,
		issuer:          cfg.Security.JWTIssuer,
	}
}

func (s *JwtServiceImpl) GenerateToken(userID uint) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiration)), // use injected expiry
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer, // use injected issuer
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

func (s *JwtServiceImpl) ValidateToken(signedToken string) (*Claims, error) {
	claims := &Claims{}
	
	token, err := jwt.ParseWithClaims(
		signedToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.secretKey, nil 
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	return claims, nil
}
