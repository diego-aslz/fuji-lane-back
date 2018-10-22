package flentities

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

// PropertyStateDraft means this property is being filled in by the user
const PropertyStateDraft = 1

var propertyStates = map[int]string{
	PropertyStateDraft: "draft",
}

// Property contains address and can have multiple units that can be booked
type Property struct {
	gorm.Model      `json:"-"`
	Name            *string    `json:"name"`
	StateID         int        `gorm:"column:state" json:"-"`
	AccountID       uint       `json:"-"`
	Account         *Account   `json:"-"`
	Address1        *string    `json:"address1"`
	Address2        *string    `json:"address2"`
	Address3        *string    `json:"address3"`
	PostalCode      *string    `json:"postalCode"`
	CityID          *int       `json:"cityID"`
	City            *City      `json:"-"`
	CountryID       *int       `json:"countryID"`
	Country         *Country   `json:"-"`
	Images          []*Image   `json:"images"`
	MinimumStay     *string    `json:"minimumStay"`
	Deposit         *string    `json:"deposit"`
	Cleaning        *string    `json:"cleaning"`
	NearestAirport  *string    `json:"nearestAirport"`
	NearestSubway   *string    `json:"nearestSubway"`
	NearbyLocations *string    `json:"nearbyLocations"`
	Overview        *string    `json:"overview"`
	Amenities       []*Amenity `json:"amenities"`
}

// State returns the state name for the property's state ID
func (p *Property) State() string {
	return propertyStates[p.StateID]
}

type propertyAlias Property

type propertyUI struct {
	ID    uint   `json:"id"`
	State string `json:"state"`
	*propertyAlias
}

func (p *propertyUI) stateID() int {
	for id, state := range propertyStates {
		if state == p.State {
			return id
		}
	}
	return 0
}

// MarshalJSON returns JSON bytes for a Property
func (p *Property) MarshalJSON() ([]byte, error) {
	return json.Marshal(propertyUI{
		ID:            p.ID,
		State:         p.State(),
		propertyAlias: (*propertyAlias)(p),
	})
}

// UnmarshalJSON loads a Property from JSON bytes
func (p *Property) UnmarshalJSON(data []byte) error {
	aux := &propertyUI{propertyAlias: (*propertyAlias)(p)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.ID = aux.ID
	p.StateID = aux.stateID()
	return nil
}
