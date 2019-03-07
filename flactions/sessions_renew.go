package flactions

import (
	"net/http"
)

// SessionsRenew gives the user a new session with a new token and expiration date
type SessionsRenew struct {
	sessionAction
}

// Perform executes the action
func (a *SessionsRenew) Perform() {
	if a.CurrentSession().RenewAfter.After(a.Now()) {
		a.Respond(http.StatusNotModified, nil)
		return
	}

	a.createSession(a.CurrentUser())
}

// NewSessionsRenew returns a new RenewSession action
func NewSessionsRenew(c Context) Action {
	return &SessionsRenew{sessionAction: sessionAction{c}}
}
