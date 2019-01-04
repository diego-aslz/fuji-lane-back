package flactions

import (
	"errors"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// BookingsCreateBody is the payload to create a booking
type BookingsCreateBody struct {
	UnitID         uint      `json:"unitID"`
	CheckInAt      time.Time `json:"checkInAt"`
	CheckOutAt     time.Time `json:"checkOutAt"`
	AdditionalInfo *string   `json:"additionalInfo"`
}

// Validate the request body
func (b *BookingsCreateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("unit", b.UnitID).Required(),
		flentities.ValidateField("check in date", b.CheckInAt).Required(),
		flentities.ValidateField("check out date", b.CheckOutAt).Required(),
	)
}

// BookingsCreate lists user bookings
type BookingsCreate struct {
	BookingsCreateBody
	Context
}

// Perform executes the action
func (a *BookingsCreate) Perform() {
	booking := &flentities.Booking{
		UserID:         a.CurrentUser().ID,
		Unit:           &flentities.Unit{ID: a.UnitID},
		UnitID:         a.UnitID,
		CheckInAt:      a.CheckInAt,
		CheckOutAt:     a.CheckOutAt,
		AdditionalInfo: a.AdditionalInfo,
	}

	errs := flentities.ValidateFields(
		flentities.ValidateField("check in date", booking.CheckInAt).After(a.Now(), "check in date should be in the future"),
		flentities.ValidateField("check out date", booking.CheckOutAt).After(booking.CheckInAt,
			"check out date should be after check in date"),
	)

	if len(errs) > 0 {
		a.RespondError(http.StatusUnprocessableEntity, errs...)
		return
	}

	if booking.CheckInAt.Before(a.Now()) {
		a.RespondError(http.StatusUnprocessableEntity, errors.New("check in date should be in the future"))
	}

	if err := a.Repository().Find(booking.Unit).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			a.RespondError(http.StatusUnprocessableEntity, errors.New("Invalid unit"))
		} else {
			a.ServerError(err)
		}

		return
	}

	if booking.Unit.PublishedAt == nil {
		a.RespondError(http.StatusUnprocessableEntity, errors.New("Invalid unit"))
		return
	}

	booking.Calculate()

	if err := a.Repository().Save(booking).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusCreated, booking)
}

// NewBookingsCreate returns a new BookingsCreate action
func NewBookingsCreate(c Context) Action {
	return &BookingsCreate{Context: c}
}
