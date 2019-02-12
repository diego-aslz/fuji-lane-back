package flservices

import (
	"fmt"
	"net/smtp"

	"github.com/nerde/fuji-lane-back/flconfig"

	"github.com/jordan-wright/email"
)

// Mailer is responsible for sending out emails
type Mailer interface {
	Send(*email.Email) error
}

// SMTPMailer sends emails through a SMTP server
type SMTPMailer struct {
	flconfig.SMTP
}

// Send an email
func (m *SMTPMailer) Send(e *email.Email) error {
	var auth smtp.Auth

	if m.Auth == "md5" {
		auth = smtp.CRAMMD5Auth(m.User, m.Password)
	} else {
		auth = smtp.PlainAuth("", m.User, m.Password, m.Host)
	}

	if e.From == "" {
		e.From = m.DefaultFrom
	}

	return e.Send(fmt.Sprintf("%s:%d", m.Host, m.Port), auth)
}

// NewSMTPMailer returns a new SMTPMailer with its configuration
func NewSMTPMailer() *SMTPMailer {
	return &SMTPMailer{flconfig.Config.SMTP}
}
