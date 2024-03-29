package flviews

import "github.com/nerde/fuji-lane-back/flentities"

// NewBookingItem returns the JSON structure for a Booking
func NewBookingItem(b *flentities.Booking) map[string]interface{} {
	return map[string]interface{}{
		"id":           b.ID,
		"propertyName": b.Unit.Property.Name,
		"unitName":     b.Unit.Name,
		"checkIn":      b.CheckIn,
		"checkOut":     b.CheckOut,
		"nights":       b.Nights,
	}
}

// NewBookingList returns an array of maps to expose a list of bookings
func NewBookingList(bs []*flentities.Booking) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, b := range bs {
		result = append(result, NewBookingItem(b))
	}

	return result
}
