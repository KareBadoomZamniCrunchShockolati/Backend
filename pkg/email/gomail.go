package email

import (
	"fmt"
	"challenge-app/internal/bootstrap"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendVerificationEmail(to, code string) error
	SendEmailChangeVerification(newEmail, code, oldEmail string) error
	SendEmailChangeConfirmation(newEmail, oldEmail string) error
}

type EmailServiceImpl struct {
	dialer *gomail.Dialer
	from   string
}

func NewEmailService(cfg *bootstrap.Env) *EmailServiceImpl {
	dialer := gomail.NewDialer(
		cfg.Email.SMTPHost,
		cfg.Email.SMTPPort,
		cfg.Email.SMTPUser,
		cfg.Email.SMTPPass,
	)

	return &EmailServiceImpl{
		dialer: dialer,
		from:   cfg.Email.From,
	}
}

func (s *EmailServiceImpl) SendVerificationEmail(to, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s\nIt expires in 5 minutes.", code))

	return s.dialer.DialAndSend(m)
}

func (s *EmailServiceImpl) SendEmailChangeVerification(newEmail, code, oldEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", newEmail)
	m.SetHeader("Subject", "Verify Your New Email Address")
	m.SetBody("text/plain",
		fmt.Sprintf("You have requested to change your email from %s to %s.\n\nYour verification code is: %s\n\nThis code expires in 10 minutes.",
			oldEmail, newEmail, code))

	return s.dialer.DialAndSend(m)
}

func (s *EmailServiceImpl) SendEmailChangeConfirmation(newEmail, oldEmail string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", newEmail)
	m.SetHeader("Subject", "Email Change Confirmed")
	m.SetBody("text/plain",
		fmt.Sprintf("Your email has been successfully changed from %s to %s.\n\nYou can now use your new email address to log in.",
			oldEmail, newEmail))

	return s.dialer.DialAndSend(m)
}
