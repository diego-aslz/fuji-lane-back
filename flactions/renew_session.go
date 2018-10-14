package flactions

import (
	"net/http"
)

// RenewSession gives the user a new session with a new token and expiration date
type RenewSession struct{}

// Perform executes the action
func (a *RenewSession) Perform(c Context) {
	if c.CurrentSession().RenewAfter.After(c.Now()) {
		c.Respond(http.StatusNotModified, nil)
		return
	}

	createSession(c, c.CurrentUser())
}
