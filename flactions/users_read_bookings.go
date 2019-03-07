package flactions

import (
	"net/http"
)

// UsersReadBookings exposes details for the logged in user
type UsersReadBookings struct {
	Context
}

// Perform executes the action
func (a *UsersReadBookings) Perform() {
	if err := a.Repository().Model(a.CurrentUser()).UpdateColumn("unread_bookings_count", 0).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, a.CurrentUser())
}

// NewUsersReadBookings returns a new UsersReadBookings action
func NewUsersReadBookings(c Context) Action {
	return &UsersReadBookings{c}
}
