package flemail

import (
	"errors"
	"net/textproto"

	"github.com/jordan-wright/email"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flutils"
)

type bookingCreated struct {
	Booking *flentities.Booking
	owner   *flentities.User
}

func (bc bookingCreated) replyTo() string {
	replyTo := bc.Booking.User.Email

	if bc.Booking.User.Name != nil && *bc.Booking.User.Name != "" {
		replyTo = *bc.Booking.User.Name + "<" + replyTo + ">"
	}

	return replyTo
}

func (bc bookingCreated) email() (*email.Email, error) {
	body, err := renderTextTemplate(bc)
	if err != nil {
		return nil, err
	}

	e := &email.Email{
		To:      []string{bc.owner.Email},
		Subject: "Booking Inquire - Fuji Lane",
		Text:    []byte(body),
		// HTML:    []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}

	e.Headers.Add("Reply-To", bc.replyTo())

	return e, nil
}

func (bc bookingCreated) User() string {
	return bc.Booking.User.Email
}

func (bc bookingCreated) Unit() string {
	propertyName := ""
	if bc.Booking.Unit.Property.Name != nil {
		propertyName = *bc.Booking.Unit.Property.Name
	}

	return propertyName + " > " + bc.Booking.Unit.Name
}

func (bc bookingCreated) PerNightPrice() string {
	return flutils.FormatCents(bc.Booking.PerNightCents)
}

func (bc bookingCreated) TotalPrice() string {
	return flutils.FormatCents(bc.Booking.TotalCents)
}

// BookingCreated returns an email to be sent when a new booking is created
func BookingCreated(b *flentities.Booking, owner *flentities.User) (*email.Email, error) {
	if b.User == nil || b.Unit == nil || b.Unit.Property == nil || owner == nil {
		return nil, errors.New("User, Unit, Property and Owner are required to send a booking email")
	}

	return bookingCreated{b, owner}.email()
}
