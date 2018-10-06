package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

func createSession(c Context, user *flentities.User) {
	createSessionWithStatus(c, user, http.StatusOK)
}

func createSessionWithStatus(c Context, user *flentities.User, status int) {
	if user.AccountID != nil && user.Account == nil {
		user.Account = &flentities.Account{}
		c.Repository().First(user.Account, *user.AccountID)
	}

	s := flentities.NewSession(user, c.Now)
	if err := s.GenerateToken(); err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(status, s)
}
