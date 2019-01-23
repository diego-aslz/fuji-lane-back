package flentities

import (
	"time"
)

// City we support
type City struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	CountryID uint    `json:"countryID"`
	Country   Country `json:"-"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// BeforeCreate to setup defaults
func (c *City) BeforeCreate() error {
	if c.Slug == "" {
		c.Slug = generateSlug(c.Name)
	}

	return nil
}
