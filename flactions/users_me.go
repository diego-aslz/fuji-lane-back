package flactions

import (
	"net/http"
)

// UsersMe exposes details for the logged in user
type UsersMe struct {
	Context
}

// Perform executes the action
func (a *UsersMe) Perform() {
	a.Respond(http.StatusOK, a.CurrentUser())
}

// NewUsersMe returns a new UsersMe action
func NewUsersMe(c Context) Action {
	return &UsersMe{c}
}
