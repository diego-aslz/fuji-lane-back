package flactions

import (
	"net/http"
	"strconv"

	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flviews"
)

// DashboardBookings lists properties' bookings
type DashboardBookings struct {
	paginatedAction
}

// Perform executes the action
func (a *DashboardBookings) Perform() {
	a.addPageDiagnostic()

	acc := a.CurrentAccount()

	bookings := []*flentities.Booking{}
	builder := a.Repository().
		Order("bookings.id desc").
		Joins("INNER JOIN units ON bookings.unit_id = units.id").
		Joins("INNER JOIN properties ON units.property_id = properties.id").
		Where("properties.account_id = ?", acc.ID).
		Preload("Unit.Property").
		Preload("User")

	if err := a.paginate(builder, a.getPage(), defaultPageSize).Find(&bookings).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Diagnostics().Add("bookings_size", strconv.Itoa(len(bookings)))

	a.Respond(http.StatusOK, flviews.NewDashboardBookingList(bookings))
}

// NewDashboardBookings returns a new DashboardBookings action
func NewDashboardBookings(c Context) Action {
	return &DashboardBookings{paginatedAction{Context: c}}
}
