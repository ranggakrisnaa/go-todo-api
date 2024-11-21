package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type MailerConfig struct {
	SmtpHost       string
	SmtpPort       int
	SenderMailName string
}

func NewMailerConfig() (*MailerConfig, error) {
	smtpHost := os.Getenv("CONFIG_SMTP_HOST")
	smtpPortStr := os.Getenv("CONFIG_SMTP_PORT")
	if smtpHost == "" && smtpPortStr == "" {
		return nil, errors.New("SMTP configuration is missing")
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return nil, err
	}

	return &MailerConfig{
		SmtpHost: smtpHost,
		SmtpPort: smtpPort,
	}, nil
}

func (c *MailerConfig) SendMail(to []string, subject, message string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", c.SenderMailName)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", "OTP Code Verification")
	m.SetBody("text/plain", message)

	d := &gomail.Dialer{Host: c.SmtpHost, Port: c.SmtpPort}

	errDial := d.DialAndSend(m)
	if errDial != nil {
		log.Fatal(errDial.Error())
	}

	return nil
}
