package email

import (
	"challenge-app/internal/domain/service"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type emailService struct {
	dialer *gomail.Dialer
	from   string
}

func NewEmailService() service.EmailService {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return &emailService{
		dialer: gomail.NewDialer(
			os.Getenv("SMTP_HOST"),
			port,
			os.Getenv("SMTP_USER"),
			os.Getenv("SMTP_PASS"),
		),
		from: os.Getenv("EMAIL_FROM"),
	}
}

func (s *emailService) SendVerificationEmail(to, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s\nIt expires in 5 minutes.", code))

	return s.dialer.DialAndSend(m)
}
