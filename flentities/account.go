package flentities

import (
	"github.com/jinzhu/gorm"
)

// Account can have several users and properties
type Account struct {
	gorm.Model
	Name      string
	Phone     *string
	Status    int
	CountryID *int
	Country   *Country
}
