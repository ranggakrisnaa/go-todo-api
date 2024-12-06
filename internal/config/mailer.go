package config

import (
	"log"

	"gopkg.in/gomail.v2"
)

type MailerConfig struct {
	SmtpHost         string
	SmtpPort         int
	SenderMailName   string
	SmtpAuthEmail    string
	SmtpAuthPassword string
}

func NewMailerConfig(cfg *MailerConfig) *MailerConfig {
	return &MailerConfig{
		SmtpHost:         cfg.SmtpHost,
		SmtpPort:         cfg.SmtpPort,
		SenderMailName:   cfg.SenderMailName,
		SmtpAuthEmail:    cfg.SmtpAuthEmail,
		SmtpAuthPassword: cfg.SmtpAuthPassword,
	}
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
