package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type MailerConfig struct {
	SmtpHost         string
	SmtpPort         int
	SenderMailName   string
	SmtpAuthEmail    string
	SmtpAuthPassword string
}

func NewMailerConfig() (*MailerConfig, error) {
	smtpHost := os.Getenv("CONFIG_SMTP_HOST")
	smtpPortStr := os.Getenv("CONFIG_SMTP_PORT")
	senderName := os.Getenv("CONFIG_SMTP_SENDER")
	smtpAuthEmail := os.Getenv("CONFIG_AUTH_EMAIL")
	smtpAuthPassword := os.Getenv("CONFIG_AUTH_PASSWORD")
	if smtpHost == "" || smtpPortStr == "" || smtpAuthEmail == "" || smtpAuthPassword == "" {
		return nil, errors.New("SMTP configuration is missing")
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return nil, err
	}

	return &MailerConfig{
		SmtpHost:         smtpHost,
		SmtpPort:         smtpPort,
		SenderMailName:   senderName,
		SmtpAuthEmail:    smtpAuthEmail,
		SmtpAuthPassword: smtpAuthPassword,
	}, nil
}

func (c *MailerConfig) SendMail(to string, subject, message string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", c.SenderMailName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", message)

	dialer := gomail.NewDialer(
		c.SmtpHost,
		c.SmtpPort,
		c.SmtpAuthEmail,
		c.SmtpAuthPassword,
	)

	errDial := dialer.DialAndSend(m)
	if errDial != nil {
		log.Fatal(errDial.Error())
	}

	return nil
}
