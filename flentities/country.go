package flentities

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

// Country we support
type Country struct {
	gorm.Model `json:"-"`
	Name       string `json:"name"`
}

type countryAlias Country

type countryUI struct {
	ID uint `json:"id"`
	*countryAlias
}

// MarshalJSON returns JSON bytes for a Country
func (c *Country) MarshalJSON() ([]byte, error) {
	return json.Marshal(countryUI{
		ID:           c.ID,
		countryAlias: (*countryAlias)(c),
	})
}

// UnmarshalJSON loads a Country from JSON bytes
func (c *Country) UnmarshalJSON(data []byte) error {
	aux := &countryUI{countryAlias: (*countryAlias)(c)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.ID = aux.ID
	return nil
}
