package model

type UserModel struct {
	ID             uint
	Username       string
	Email          string
	PasswordHash   string
	Bio            string
	ProfilePicture string
	Verified       bool
}
