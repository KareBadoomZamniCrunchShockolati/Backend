package dto

// InitiateEmailChangeRequest for starting email change process
type InitiateEmailChangeRequest struct {
	NewEmail string `json:"new_email" binding:"required,email"`
}

// VerifyEmailChangeRequest for completing email change
type VerifyEmailChangeRequest struct {
	OldEmail string `json:"old_email" binding:"required,email"`
	NewEmail string `json:"new_email" binding:"required,email"`
	Code     string `json:"code" binding:"required,len=6"`
}

// EmailChangeResponse for email change responses
type EmailChangeResponse struct {
	Message string `json:"message"`
	Email   string `json:"email,omitempty"`
}
