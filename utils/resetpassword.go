package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
)

func SendResetPasswordEmail(to string, token string) error {

	email := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	if to == "" {
		return fmt.Errorf("email address is empty")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Reset Your Password")
	resetLink := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", token)
	m.SetBody("text/html", fmt.Sprintf(`
		<h2>Reset your password</h2>
		<p>Click the link below to reset your password:</p>
		<a href="%s">%s</a>
	`, resetLink, resetLink))

	d := gomail.NewDialer("smtp.gmail.com", 587, email, password)

	return d.DialAndSend(m)
}
