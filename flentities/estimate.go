package flentities

import "math"

// Estimate represent a price estimate to book a unit for a given amount of nights
type Estimate struct {
	Nights          int
	NightPriceCents int
	TotalCents      int
}

// NewEstimate estimates the price to book the given unit for the given nights
func NewEstimate(unit *Unit, nights int) Estimate {
	p := unit.PriceFor(nights)

	return Estimate{
		Nights:          nights,
		NightPriceCents: (int)(math.Round(p.PerNightCents())),
		TotalCents:      int(math.Round(float64(nights) * p.PerNightCents())),
	}
}
