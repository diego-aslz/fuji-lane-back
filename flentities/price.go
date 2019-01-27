package flentities

import "time"

// PerNightPriceSQL is the expression we use to calculate the Per Night for a Price in the database
const PerNightPriceSQL = "prices.cents / prices.min_nights"

// Price represents a unit price under specific conditions
type Price struct {
	ID        uint      `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	UnitID    uint      `json:"-"`
	Unit      *Unit     `json:"-"`
	MinNights int       `json:"minNights"`
	Cents     int       `json:"cents"`
}

// PerNightCents returns the price per night for this price definition
func (p Price) PerNightCents() float64 {
	return float64(p.Cents) / float64(p.MinNights)
}
