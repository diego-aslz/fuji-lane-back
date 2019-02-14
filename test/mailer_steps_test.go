package fujilane

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jordan-wright/email"
)

var mailer *fakeMailer

type fakeMailer struct {
	emails []*email.Email
}

func (m *fakeMailer) Send(e *email.Email) error {
	m.emails = append(m.emails, e)
	return nil
}

func assertEmail(to string, doc *gherkin.DocString) error {
	for _, e := range mailer.emails {
		if strings.Join(e.To, ", ") != to {
			continue
		}

		body := strings.TrimSpace(string(e.Text))
		if body == doc.Content {
			return nil
		}

		return fmt.Errorf("Expected email to have body:\n%s\nBut got body:\n%s", doc.Content, body)
	}

	return errors.New("No emails matched the criteria")
}

func assertNoEmails() error {
	count := len(mailer.emails)
	if count == 0 {
		return nil
	}

	return fmt.Errorf("Expected 0 emails to have been sent, but actually sent %d", count)
}

func MailerContext(s *godog.Suite) {
	s.BeforeScenario(func(_ interface{}) {
		mailer = &fakeMailer{}
		application.Mailer = mailer
	})

	s.Step(`^"([^"]*)" should have received the following email:$`, assertEmail)
	s.Step(`^no emails should have been sent$`, assertNoEmails)
}
