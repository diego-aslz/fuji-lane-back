package flentities

import (
	"github.com/jinzhu/gorm"
)

// Account can have several users and properties
type Account struct {
	gorm.Model `json:"-"`
	Name       string   `json:"name"`
	Phone      *string  `json:"phone"`
	Status     int      `json:"-"`
	CountryID  *int     `json:"country_id"`
	Country    *Country `json:"-"`
}
