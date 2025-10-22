package security

import (
	"golang.org/x/crypto/bcrypt"
)


type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(hash, password string) bool
}

type PasswordServiceImpl struct{}

func NewPasswordService() *PasswordServiceImpl {
	return &PasswordServiceImpl{}
}

// HashPassword generates a bcrypt hash of the password.
func (s *PasswordServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plaintext password with a hashed password.
func (s *PasswordServiceImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
