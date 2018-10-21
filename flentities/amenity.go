package flentities

import "time"

// AmenityType represents a recognized amenity type
type AmenityType struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	ForProperties bool   `json:"-"`
	ForUnits      bool   `json:"-"`
}

// AmenityTypes supported by the system
var AmenityTypes = []*AmenityType{
	{"daycare", "Daycare", true, false},
	{"gym", "Gym", true, false},
	{"meeting_rooms", "Meeting Rooms", true, false},
	{"pool", "Pool", true, false},
	{"restaurant", "Restaurant", true, false},

	{"air_conditioning", "Air Conditioning", false, true},
	{"bathrobes", "Bathrobes", false, true},
	{"blackout_curtains", "Blackout Curtains", false, true},
	{"housekeeping", "Daily Housekeeping", false, true},
	{"desk", "Desk", false, true},
	{"dvd", "DVD Player", false, true},
	{"minibar", "Minibar", false, true},
	{"phone", "Phone", false, true},
	{"toilet", "Toilet", false, true},
}

// Amenity associates a Property or Unit to an amenity type
type Amenity struct {
	ID         uint        `json:"id"`
	CreatedAt  time.Time   `json:"-"`
	UpdatedAt  time.Time   `json:"-"`
	PropertyID uint        `json:"propertyID"`
	Property   Property    `json:"-"`
	Type       AmenityType `json:"type"`
}
