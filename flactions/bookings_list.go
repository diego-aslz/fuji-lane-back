package flactions

import (
	"net/http"
	"strconv"

	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flviews"
)

const bookingsPageSize = defaultPageSize

// BookingsList lists user bookings
type BookingsList struct {
	paginatedAction
}

// Perform executes the action
func (a *BookingsList) Perform() {
	user := a.CurrentUser()

	bookings := []*flentities.Booking{}
	err := a.paginate(a.Repository().Order("check_in_at desc").Preload("Unit"), a.page(), bookingsPageSize).Find(
		&bookings, map[string]interface{}{"user_id": user.ID}).Error
	if err != nil {
		a.ServerError(err)
		return
	}

	a.Diagnostics().Add("bookings_size", strconv.Itoa(len(bookings)))

	a.Respond(http.StatusOK, flviews.NewBookingList(bookings))
}

// NewBookingsList returns a new BookingsList action
func NewBookingsList(c Context) Action {
	return &BookingsList{paginatedAction: paginatedAction{c}}
}
