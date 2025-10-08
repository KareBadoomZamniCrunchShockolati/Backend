package dto

import "github.com/google/uuid"


// UpdateProfileRequest is used to update user details.
type UpdateProfileRequest struct {
	Username string `json:"username,omitempty"`
	Bio      string `json:"bio,omitempty"`
	NewEmail string `json:"new_email,omitempty"`

	// Note: Email changes should typically be handled via a separate, verified process.
}

// UserResponse is used for sending profile data back to the client.
type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Bio      string    `json:"bio"`
}
