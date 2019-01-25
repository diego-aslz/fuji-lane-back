package flentities

import (
	"math"
	"time"
)

// Booking represents a request from a User to book a Unit for a period of time
type Booking struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"-"`
	UserID          uint      `json:"userID"`
	User            *User     `json:"-"`
	UnitID          uint      `json:"unitID"`
	Unit            *Unit     `json:"-"`
	CheckIn         Date      `json:"checkIn"`
	CheckOut        Date      `json:"checkOut"`
	Message         *string   `json:"message"`
	NightPriceCents int       `json:"nightPriceCents"`
	Nights          int       `json:"nights"`
	ServiceFeeCents int       `json:"serviceFeeCents"`
	TotalCents      int       `json:"totalCents"`
}

// Calculate fills in calculated fields, like prices
func (b *Booking) Calculate() {
	b.calculateNights()
	b.calculatePrice()
}

func (b *Booking) calculateNights() {
	b.Nights = b.CheckOut.NightsFrom(b.CheckIn)
}

func (b *Booking) calculatePrice() {
	if b.Unit == nil {
		return
	}

	p := b.Unit.PriceFor(b.Nights)

	if p.Cents == 0 {
		return
	}

	b.NightPriceCents = (int)(math.Round(p.PerNightCents()))
	b.TotalCents = (int)(math.Round((float64)(b.Nights) * p.PerNightCents()))
}
