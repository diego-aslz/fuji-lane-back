package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type bookingRow struct {
	flentities.Booking
	Unit string
	User string
}

func tableRowToBooking(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*bookingRow)

	if row.User != "" {
		row.Booking.User = &flentities.User{}
		findBy(row.Booking.User, "email", row.User)
	}

	return &row.Booking, loadAssociationByName(&row.Booking, "Unit", row.Unit)
}

func BookingContext(s *godog.Suite) {
	s.Step(`^the following bookings:$`, createFromTableStep(new(bookingRow), tableRowToBooking))
	s.Step(`^I get dashboard details for:$`, performGETWithParamsStep(flweb.DashboardPath))
}
