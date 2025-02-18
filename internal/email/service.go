package email

import (
	"math/rand"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Service struct {
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
}

func NewEmailService(host string, port int, username, password string) *Service {
	return &Service{
		smtpHost:     host,
		smtpPort:     port,
		smtpUsername: username,
		smtpPassword: password,
	}
}

func (s *Service) SendVerificationCode(to string) (string, error) {
	code := generateVerificationCode()

	m := gomail.NewMessage()
	m.SetHeader("From", s.smtpUsername)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "验证码")
	m.SetBody("text/html", "您的验证码是："+code)

	d := gomail.NewDialer(s.smtpHost, s.smtpPort, s.smtpUsername, s.smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return "", err
	}

	return code, nil
}

func generateVerificationCode() string {
	code := rand.Intn(999999)
	return strconv.Itoa(code)
} 