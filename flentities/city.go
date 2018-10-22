package flentities

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

// City we support
type City struct {
	gorm.Model `json:"-"`
	Name       string  `json:"name"`
	CountryID  uint    `json:"countryID"`
	Country    Country `json:"-"`
}

type cityAlias City

type cityUI struct {
	ID uint `json:"id"`
	*cityAlias
}

// MarshalJSON returns JSON bytes for a City
func (c *City) MarshalJSON() ([]byte, error) {
	return json.Marshal(cityUI{
		ID:        c.ID,
		cityAlias: (*cityAlias)(c),
	})
}

// UnmarshalJSON loads a City from JSON bytes
func (c *City) UnmarshalJSON(data []byte) error {
	aux := &cityUI{cityAlias: (*cityAlias)(c)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.ID = aux.ID
	return nil
}
