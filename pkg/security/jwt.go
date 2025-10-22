package security

import (
	"fmt"
	"time"
	"challenge-app/internal/bootstrap"
	"github.com/golang-jwt/jwt/v5"
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
}

// NewJWTService creates a new JWTService implementation.
func NewJWTService(secretKey string, tokenExpiration time.Duration) *JwtServiceImpl {
	return &JwtServiceImpl{
		secretKey:       []byte(secretKey),
		tokenExpiration: tokenExpiration,
	}
}

func (s *JwtServiceImpl) GenerateToken(userID uint) (string, error) {

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(bootstrap.JWTTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    bootstrap.JWTIssuer,
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
