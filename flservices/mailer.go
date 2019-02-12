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
	return e.Send(fmt.Sprintf("%s:%d", m.Host, m.Port), smtp.CRAMMD5Auth(m.User, m.Password))
}

// NewSMTPMailer returns a new SMTPMailer with its configuration
func NewSMTPMailer() *SMTPMailer {
	return &SMTPMailer{flconfig.Config.SMTP}
}
