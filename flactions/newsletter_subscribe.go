package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flservices"
)

// NewsletterSubscribe subscribes an email to our Newsletter
type NewsletterSubscribe struct {
	Context `json:"-"`
	flservices.SendgridContact
	Sendgrid flservices.Sendgrid `json:"-"`
}

// Validate the request body
func (a *NewsletterSubscribe) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("email", a.Email).Required().Email(),
	)
}

// Perform executes the action
func (a *NewsletterSubscribe) Perform() {
	if err := a.Sendgrid.SubscribeNewsletter(a.SendgridContact); err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusCreated, nil)
}

// NewNewsletterSubscribe returns a new NewsletterSubscribe action
func NewNewsletterSubscribe(c Context, s flservices.Sendgrid) Action {
	return &NewsletterSubscribe{Context: c, Sendgrid: s}
}
