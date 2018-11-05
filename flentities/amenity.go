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
	ID         uint      `json:"-"`
	CreatedAt  time.Time `json:"-"`
	PropertyID *uint     `json:"-"`
	Property   *Property `json:"-"`
	UnitID     *uint     `json:"-"`
	Unit       *Unit     `json:"-"`
	Type       string    `json:"type"`
	Name       *string   `json:"name"`
}

// IsValidAmenity checks if the given amenity type and name are to be accepted by the system. Custom amenities are
// defined by the user and need a name. Other amenities need to be recognized by the system, that is, they need to
// be in AmenityTypes slice
func IsValidAmenity(aType, name string) bool {
	if aType == "custom" {
		return name != ""
	}

	for _, typ := range AmenityTypes {
		if typ.Code == aType {
			return true
		}
	}

	return false
}
