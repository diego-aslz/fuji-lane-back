package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flactions"
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

func performGETStepWithPage(path string) func(string) error {
	return func(page string) error {
		return perform("GET", path+"?page="+page, nil)
	}
}

type bookingBodyTable struct {
	flactions.BookingsCreateBody
	Unit string
}

func requestBookingsCreate(table *gherkin.DataTable) error {
	bbt, err := assist.CreateInstance(new(bookingBodyTable), table)
	if err != nil {
		return err
	}

	body := bbt.(*bookingBodyTable)

	if body.Unit != "" {
		unit := &flentities.Unit{}
		if err := findByName(unit, body.Unit); err != nil {
			return err
		}

		body.UnitID = unit.ID
	}

	return performPOST(flweb.BookingsPath, body)
}

func bookingToTableRow(r *flentities.Repository, a interface{}) (interface{}, error) {
	b := a.(*flentities.Booking)

	b.Unit = &flentities.Unit{}
	err := r.Model(b).Association("Unit").Find(b.Unit).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	b.User = &flentities.User{}
	err = r.Model(b).Association("User").Find(b.User).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return &bookingRow{Booking: *b, Unit: b.Unit.Name, User: b.User.Email}, err
}

func BookingContext(s *godog.Suite) {
	s.Step(`^the following bookings:$`, createFromTableStep(new(bookingRow), tableRowToBooking))
	s.Step(`^I get dashboard details for:$`, performGETWithParamsStep(flweb.DashboardPath))
	s.Step(`^I list my bookings$`, performGETStep(flweb.BookingsPath))
	s.Step(`^I list my bookings for page "([^"]*)"$`, performGETStepWithPage(flweb.BookingsPath))
	s.Step(`^I submit the following booking:$`, requestBookingsCreate)
	s.Step(`^I should have the following bookings:$`, assertDatabaseRecordsStep(&[]*flentities.Booking{},
		bookingToTableRow))
}
