package flemail

import (
	"fmt"
	"net/textproto"
	"strconv"

	"github.com/jordan-wright/email"
	"github.com/nerde/fuji-lane-back/flentities"
)

type bookingCreated struct {
	booking *flentities.Booking
	owner   *flentities.User
}

func (bc bookingCreated) replyTo() string {
	replyTo := bc.booking.User.Email

	if bc.booking.User.Name != nil && *bc.booking.User.Name != "" {
		replyTo = *bc.booking.User.Name + "<" + replyTo + ">"
	}

	return replyTo
}

func (bc bookingCreated) email() *email.Email {
	e := &email.Email{
		To:      []string{bc.owner.Email},
		Subject: "[Fuji Lane] Inquire",
		Text:    []byte(bc.textBody()),
		// HTML:    []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}

	e.Headers.Add("Reply-To", bc.replyTo())

	return e
}

func (bc bookingCreated) textBody() string {
	body := "Hi there,\n\nYou received a new booking request:\n\n"

	body += "* User: " + bc.user() + "\n"
	body += "* Unit: " + bc.unit() + "\n"
	body += "* Check In: " + bc.booking.CheckIn.String() + "\n"
	body += "* Check Out: " + bc.booking.CheckOut.String() + "\n"
	body += "* Nights: " + strconv.Itoa(bc.booking.Nights) + "\n"
	body += "* Price: $" + fmt.Sprint(float32(bc.booking.PerNightCents)/100.0) + "/night\n"
	body += "* Total: $" + fmt.Sprint(float32(bc.booking.TotalCents)/100.0) + "\n"

	body += "\nRespond to this email to get in touch with them."

	return body
}

func (bc bookingCreated) user() string {
	return bc.booking.User.Email
}

func (bc bookingCreated) unit() string {
	propertyName := ""
	if bc.booking.Unit.Property.Name != nil {
		propertyName = *bc.booking.Unit.Property.Name
	}

	return propertyName + " > " + bc.booking.Unit.Name
}

// BookingCreated returns an email to be sent when a new booking is created
func BookingCreated(b *flentities.Booking, owner *flentities.User) *email.Email {
	if b.User == nil || b.Unit == nil || b.Unit.Property == nil || owner == nil {
		return nil
	}

	return bookingCreated{b, owner}.email()
}
