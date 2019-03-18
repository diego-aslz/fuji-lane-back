package main

import (
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flemail"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flservices"
)

func main() {
	flconfig.LoadConfiguration()

	uName := "Requestor"
	checkIn, _ := flentities.ParseDate("2019-02-01")
	checkOut, _ := flentities.ParseDate("2019-02-08")
	b := &flentities.Booking{
		User: &flentities.User{
			Name:  &uName,
			Email: "requestor@fujilane.com",
		},
		Unit: &flentities.Unit{
			Name:     "Standard Apartment",
			Property: &flentities.Property{Name: "Awesome Property"},
		},
		CheckIn:       checkIn,
		CheckOut:      checkOut,
		Nights:        7,
		PerNightCents: 10000,
		TotalCents:    70000,
	}

	owner := &flentities.User{
		Email: "owner@fujilane.com",
	}

	e, er := flemail.BookingCreated(b, owner)
	if er != nil {
		panic(er)
	}

	if err := flservices.NewSMTPMailer().Send(e); err != nil {
		panic(err)
	}
}
