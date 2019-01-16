package flactions

import (
	"net/http"
	"strconv"

	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flviews"
)

// BookingsList lists user bookings
type BookingsList struct {
	paginatedAction
}

// Perform executes the action
func (a *BookingsList) Perform() {
	a.addPageDiagnostic()

	user := a.CurrentUser()

	bookings := []*flentities.Booking{}
	err := a.paginate(a.Repository().Order("check_in desc").Preload("Unit.Property"), a.getPage(), defaultPageSize).Find(
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
	return &BookingsList{paginatedAction: paginatedAction{Context: c}}
}
