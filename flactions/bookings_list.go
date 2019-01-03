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
func (a *BookingsList) Perform(c Context) {
	user := c.CurrentUser()

	bookings := []*flentities.Booking{}
	err := a.paginate(c.Repository().Order("check_in_at desc").Preload("Unit"), a.page(c), bookingsPageSize).Find(
		&bookings, map[string]interface{}{"user_id": user.ID}).Error
	if err != nil {
		c.ServerError(err)
		return
	}

	c.Diagnostics().Add("bookings_size", strconv.Itoa(len(bookings)))

	c.Respond(http.StatusOK, flviews.NewBookingList(bookings))
}
