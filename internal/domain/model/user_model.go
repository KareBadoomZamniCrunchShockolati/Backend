package model

// User is the core entity.
type UserModel struct {
	ID             uint
	Username       string
	Email          string
	PasswordHash   string
	Bio            string
	ProfilePicture string
	Verified       bool
}
