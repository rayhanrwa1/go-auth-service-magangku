package service

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailService struct {
	SMTPHost string
	SMTPPort string
	From     string
	Password string
}

func NewEmailService() *EmailService {
	return &EmailService{
		SMTPHost: os.Getenv("SMTP_HOST"),
		SMTPPort: os.Getenv("SMTP_PORT"),
		From:     os.Getenv("SMTP_FROM"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

func (s *EmailService) SendResetPasswordEmail(to, token string) error {

	// TODO : Link Password disini di generate
	resetLink := fmt.Sprintf("https://www.magangku.web.id/lupapassword?oobCode=%s", token)
	
	subject := "Reset Your Password - Magangku"
	body := fmt.Sprintf(`
		<h2>Reset Password Request</h2>
		<p>Halo,</p>
		<p>Anda telah meminta untuk mereset password akun Magangku Anda.</p>
		<p>Klik tombol di bawah ini untuk mereset password:</p>
		<a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;">Reset Password</a>
		<p>Atau copy link berikut ke browser Anda:</p>
		<p>%s</p>
		<p>Link ini akan kadaluarsa dalam 1 jam.</p>
		<p>Jika Anda tidak meminta reset password, abaikan email ini.</p>
		<br>
		<p>Salam,<br>Tim Magangku</p>
	`, resetLink, resetLink)

	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-version: 1.0;\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
			"\r\n"+
			"%s\r\n",
		s.From, to, subject, body,
	))

	auth := smtp.PlainAuth("", s.From, s.Password, s.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)

	return smtp.SendMail(addr, auth, s.From, []string{to}, message)
}