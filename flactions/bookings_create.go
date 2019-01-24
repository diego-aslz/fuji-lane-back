package flactions

import (
	"errors"
	"net/http"

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
	booking := &flentities.Booking{
		UserID:   a.CurrentUser().ID,
		Unit:     &flentities.Unit{ID: a.UnitID},
		UnitID:   a.UnitID,
		CheckIn:  a.CheckIn,
		CheckOut: a.CheckOut,
	}

	if a.Message != nil && *a.Message != "" {
		booking.Message = a.Message
	}

	if booking.CheckIn.Before(a.Now()) {
		a.RespondError(http.StatusUnprocessableEntity, errors.New("check in date should be in the future"))
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

	if err := a.Repository().Save(booking).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusCreated, booking)
}

func (a *BookingsCreate) invalidUnit() {
	a.RespondError(http.StatusUnprocessableEntity, errors.New("unit is invalid"))
}

// NewBookingsCreate returns a new BookingsCreate action
func NewBookingsCreate(c Context) Action {
	return &BookingsCreate{Context: c}
}
