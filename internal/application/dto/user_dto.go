package dto

// UpdateProfileRequest is used to update user details.
type UpdateProfileRequest struct {
	Username       string `json:"username,omitempty"`
	Bio            string `json:"bio,omitempty"`
	NewEmail       string `json:"new_email,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

// UserResponse is used for sending profile data back to the client.
type UserResponse struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Bio            string `json:"bio"`
	ProfilePicture string `json:"profile_picture"`
}

type UserPreviewDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

