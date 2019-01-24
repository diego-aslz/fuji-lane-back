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
	b.Nights = int(math.Ceil(b.CheckOut.Sub(b.CheckIn.Time).Hours() / 24))
}

func (b *Booking) calculatePrice() {
	if b.Unit == nil {
		return
	}

	perNight := (float64)(*b.Unit.BasePriceCents)
	switch {
	case b.Nights >= 365:
		perNight = b.safePrice(b.Unit.TwelveMonthsPriceCents, 365.0)
		break
	case b.Nights >= 180:
		perNight = b.safePrice(b.Unit.SixMonthsPriceCents, 180.0)
		break
	case b.Nights >= 90:
		perNight = b.safePrice(b.Unit.ThreeMonthsPriceCents, 90.0)
		break
	case b.Nights >= 7:
		perNight = b.safePrice(b.Unit.OneWeekPriceCents, 7.0)
		break
	case b.Nights == 1:
		perNight = b.safePrice(b.Unit.OneNightPriceCents, 1.0)
		break
	}

	b.NightPriceCents = (int)(math.Round(perNight))
	b.TotalCents = (int)(math.Round((float64)(b.Nights) * perNight))
}

func (b *Booking) safePrice(cents *int, days float64) float64 {
	if cents == nil || *cents == 0 {
		return (float64)(*b.Unit.BasePriceCents)
	}

	return (float64)(*cents) / days
}
