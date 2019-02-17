package flviews

import "github.com/nerde/fuji-lane-back/flentities"

// NewDashboardBookingItem returns the JSON structure for a Booking for the Dashboard
func NewDashboardBookingItem(b *flentities.Booking) map[string]interface{} {
	bm := NewBookingItem(b)
	bm["perNightCents"] = b.PerNightCents
	bm["totalCents"] = b.TotalCents
	bm["user"] = map[string]interface{}{
		"name":  b.User.Name,
		"email": b.User.Email,
	}

	return bm
}

// NewDashboardBookingList returns an array of maps to expose a list of bookings for the Dashboard
func NewDashboardBookingList(bs []*flentities.Booking) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, b := range bs {
		result = append(result, NewDashboardBookingItem(b))
	}

	return result
}
