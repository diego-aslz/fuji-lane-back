package flentities

import (
	"encoding/json"
	"time"
)

// Amenity associates a Property or Unit to an amenity type
type Amenity struct {
	ID         uint      `json:"id"`
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

// MarshalJSON to create a JSON with calculated name
func (a Amenity) MarshalJSON() ([]byte, error) {
	str := AmenityName(a)
	a.Name = &str

	return json.Marshal(&struct {
		ID   uint   `json:"id"`
		Type string `json:"type"`
		Name string `json:"name"`
	}{a.ID, a.Type, *a.Name})
}

// AmenityName returns the calculated amenity's name. It returns the persisted name if it's a custom amenity
func AmenityName(a Amenity) string {
	if a.Type != "custom" && a.Name == nil {
		for _, typ := range AmenityTypes {
			if typ.Code == a.Type {
				return typ.Name
			}
		}
	}

	if a.Name == nil {
		return ""
	}

	return *a.Name
}
