package security

import "golang.org/x/crypto/bcrypt"

type PasswordService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type passwordServiceImpl struct{}

func NewPasswordService() PasswordService {
	return &passwordServiceImpl{}
}

// HashPassword generates a bcrypt hash of the password.
func (s *passwordServiceImpl) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plaintext password with a hashed password.
func (s *passwordServiceImpl) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
