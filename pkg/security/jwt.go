package security

import (
	"fmt"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
	// Secret key From environment variable
	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtSecret) == 0 {
		return "", fmt.Errorf("JWT_SECRET_KEY is not set")
	}

	// Set expiration time
	expirationTime := time.Now().Add(time.Hour * 24) 

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "challenge-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}


func ValidateToken(tokenString string) (uint, error) {
	jwtSecret := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(jwtSecret) == 0 {
		return 0, fmt.Errorf("JWT_SECRET_KEY is not set")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("token validation failed: %w", err)
	}

	return claims.UserID, nil
}