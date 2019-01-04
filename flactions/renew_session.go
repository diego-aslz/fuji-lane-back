package flactions

import (
	"net/http"
)

// RenewSession gives the user a new session with a new token and expiration date
type RenewSession struct {
	sessionAction
}

// Perform executes the action
func (a *RenewSession) Perform() {
	if a.CurrentSession().RenewAfter.After(a.Now()) {
		a.Respond(http.StatusNotModified, nil)
		return
	}

	a.createSession(a.CurrentUser())
}

// NewRenewSession returns a new RenewSession action
func NewRenewSession(c Context) Action {
	return &RenewSession{sessionAction: sessionAction{c}}
}
