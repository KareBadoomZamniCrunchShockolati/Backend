package security

import (
	"fmt"
	"time"

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
type jwtServiceImpl struct {
	secretKey       []byte
	tokenExpiration time.Duration
}

// NewJWTService creates a new JWTService implementation.
func NewJWTService(secretKey string, tokenExpiration time.Duration) JWTService {
	return &jwtServiceImpl{
		secretKey:       []byte(secretKey),
		tokenExpiration: tokenExpiration,
	}
}

func (s *jwtServiceImpl) GenerateToken(userID uint) (string, error) {
	expirationTime := time.Now().Add(s.tokenExpiration)

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "challenge-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey) 
}

func (s *jwtServiceImpl) ValidateToken(signedToken string) (*Claims, error) {
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
