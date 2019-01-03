package flentities

import (
	"math"
	"time"
)

// Booking represents a request from a User to book a Unit for a period of time
type Booking struct {
	ID              uint      `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	UserID          uint      `json:"userID"`
	User            *User     `json:"-"`
	UnitID          uint      `json:"unitID"`
	Unit            *Unit     `json:"-"`
	CheckInAt       time.Time `json:"checkInAt"`
	CheckOutAt      time.Time `json:"checkOutAt"`
	AdditionalInfo  *string   `json:"additionalInfo"`
	NightPriceCents int       `json:"nightPriceCents"`
	Nights          int       `json:"nights"`
	ServiceFeeCents int       `json:"serviceFeeCents"`
	TotalCents      int       `json:"totalCents"`
}

// Calculate fills in calculated fields, like prices
func (b *Booking) Calculate() {
	b.Nights = int(math.Ceil(b.CheckOutAt.Sub(b.CheckInAt).Hours() / 24))

	if b.Unit == nil || b.Unit.BasePriceCents == nil {
		return
	}

	b.NightPriceCents = *b.Unit.BasePriceCents
	b.TotalCents = b.NightPriceCents * b.Nights
}
