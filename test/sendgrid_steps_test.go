package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/nerde/fuji-lane-back/flservices"
)

var sendgrid *fakeSendgrid

type fakeSendgrid struct {
	subscriptions []*flservices.SendgridContact
}

func (s *fakeSendgrid) SubscribeNewsletter(c flservices.SendgridContact) error {
	s.subscriptions = append(s.subscriptions, &c)
	return nil
}

func SendgridContext(s *godog.Suite) {
	s.BeforeScenario(func(_ interface{}) {
		sendgrid = &fakeSendgrid{}
		application.Sendgrid = sendgrid
	})
}
