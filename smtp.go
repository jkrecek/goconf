package goconf

import (
	"fmt"
	"net/smtp"
)

type SmtpConnection struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func (sc *SmtpConnection) GetConnectionString() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}

type EmailMessage interface {
	GetFrom() string
	GetTo() []string
	GetBytes() []byte
}

func (sc *SmtpConnection) SendEmail(em EmailMessage) error {
	auth := smtp.Auth(nil)
	if len(sc.Username) > 0 {
		auth = smtp.PlainAuth("", sc.Username, sc.Password, sc.Host)
	}

	err := smtp.SendMail(sc.GetConnectionString(), auth, em.GetFrom(), em.GetTo(), em.GetBytes())
	return err
}
