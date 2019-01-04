package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

type sessionAction struct {
	Context
}

func (a sessionAction) createSession(user *flentities.User) {
	a.createSessionWithStatus(user, http.StatusOK)
}

func (a sessionAction) createSessionWithStatus(user *flentities.User, status int) {
	if user.AccountID != nil && user.Account == nil {
		user.Account = &flentities.Account{}
		a.Repository().First(user.Account, *user.AccountID)
	}

	s := flentities.NewSession(user, a.Now)
	if err := s.GenerateToken(); err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(status, s)
}
