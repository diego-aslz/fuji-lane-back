package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"

	"github.com/nerde/fuji-lane-back/flviews"
)

// ProfileShow exposes details for the logged in user
type ProfileShow struct {
	Context
}

// Perform executes the action
func (a *ProfileShow) Perform() {
	user := a.CurrentUser()

	if user.AccountID != nil {
		user.Account = &flentities.Account{}
		user.Account.ID = *user.AccountID
		if err := a.Repository().Find(user.Account).Error; err != nil {
			a.ServerError(err)
			return
		}
	}

	a.Respond(http.StatusOK, flviews.NewProfile(user))
}

// NewProfileShow returns a new UsersMe action
func NewProfileShow(c Context) Action {
	return &ProfileShow{c}
}
