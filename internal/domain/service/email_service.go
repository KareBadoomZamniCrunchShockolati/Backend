package service

type EmailService interface {
	SendVerificationEmail(to, code string) error
	SendEmailChangeVerification(newEmail, code, oldEmail string) error
	SendEmailChangeConfirmation(newEmail, oldEmail string) error
}
