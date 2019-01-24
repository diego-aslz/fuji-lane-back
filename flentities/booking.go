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

	lastMinNights := 0
	perNight := 0.0

	for _, p := range b.Unit.Prices {
		if p.MinNights > lastMinNights && p.MinNights <= b.Nights {
			lastMinNights = p.MinNights
			perNight = float64(p.Cents) / float64(p.MinNights)
		}
	}

	b.NightPriceCents = (int)(math.Round(perNight))
	b.TotalCents = (int)(math.Round((float64)(b.Nights) * perNight))
}
