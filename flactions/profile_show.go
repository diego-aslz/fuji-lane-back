package flactions

import (
	"net/http"
)

// ProfileShow exposes details for the logged in user
type ProfileShow struct {
	Context
}

// Perform executes the action
func (a *ProfileShow) Perform() {
	a.Respond(http.StatusOK, a.CurrentUser())
}

// NewProfileShow returns a new UsersMe action
func NewProfileShow(c Context) Action {
	return &ProfileShow{c}
}
