package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/fljobs"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// BookingsCreateBody is the payload to create a booking
type BookingsCreateBody struct {
	UnitID   uint            `json:"unitID"`
	CheckIn  flentities.Date `json:"checkIn"`
	CheckOut flentities.Date `json:"checkOut"`
	Message  *string         `json:"message"`
}

// BookingsCreate lists user bookings
type BookingsCreate struct {
	BookingsCreateBody
	Context
	jobs *fljobs.Application
}

// Validate the request body
func (a *BookingsCreate) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("unit", a.UnitID).Required(),
		flentities.ValidateField("check in", a.CheckIn.Time).Required().After(a.Now(), "check in should be in the future"),
		flentities.ValidateField("check out", a.CheckOut.Time).Required().After(a.CheckIn.Time,
			"check out should be after check in"),
	)
}

// Perform executes the action
func (a *BookingsCreate) Perform() {
	unit := &flentities.Unit{ID: a.UnitID}

	if err := a.Repository().Preload("Property").Find(unit).Error; err != nil {
		a.ServerError(err)
		return
	}

	booking := &flentities.Booking{
		User:     a.CurrentUser(),
		UserID:   a.CurrentUser().ID,
		Unit:     unit,
		UnitID:   unit.ID,
		CheckIn:  a.CheckIn,
		CheckOut: a.CheckOut,
	}

	if a.Message != nil && *a.Message != "" {
		booking.Message = a.Message
	}

	if booking.CheckIn.Before(a.Now()) {
		a.RespondError(http.StatusUnprocessableEntity, errors.New("check in date should be in the future"))
		return
	}

	if err := a.Repository().Preload("Prices").Find(booking.Unit).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			a.invalidUnit()
		} else {
			a.ServerError(err)
		}

		return
	}

	if booking.Unit.PublishedAt == nil {
		a.invalidUnit()
		return
	}

	booking.Calculate()

	a.Repository().Transaction(func(tx *flentities.Repository) {
		if err := a.Repository().Save(booking).Error; err != nil {
			a.ServerError(err)
			tx.Rollback()
			return
		}

		err := a.Repository().
			Table("users").
			Where("account_id = ?", unit.Property.AccountID).
			UpdateColumn("unread_bookings_count", gorm.Expr("unread_bookings_count + ?", 1)).
			Error

		if err != nil {
			a.ServerError(err)
			tx.Rollback()
			return
		}

		err = a.Repository().
			Table("accounts").
			Where("id = ?", unit.Property.AccountID).
			UpdateColumn("bookings_count", gorm.Expr("bookings_count + ?", 1)).
			Error

		if err != nil {
			a.ServerError(err)
			tx.Rollback()
			return
		}

		var jobID string
		if jobID, err = a.jobs.EnqueueBookingCreated(booking.ID); err != nil {
			a.ServerError(err)
			tx.Rollback()
			return
		}

		a.Diagnostics().Add("BookingCreatedJobID", jobID)

		if err = tx.Commit().Error; err != nil {
			a.ServerError(err)
			return
		}

		a.Respond(http.StatusCreated, booking)
	})
}

func (a *BookingsCreate) invalidUnit() {
	a.RespondError(http.StatusUnprocessableEntity, errors.New("unit is invalid"))
}

// NewBookingsCreate returns a new BookingsCreate action
func NewBookingsCreate(c Context, jobs *fljobs.Application) Action {
	return &BookingsCreate{Context: c, jobs: jobs}
}
