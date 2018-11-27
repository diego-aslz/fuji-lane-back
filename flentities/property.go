package flentities

import (
	"time"
)

// Property contains address and can have multiple units that can be booked
type Property struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `json:"-"`
	UpdatedAt       time.Time  `json:"-"`
	DeletedAt       *time.Time `json:"-"`
	Name            *string    `json:"name"`
	PublishedAt     *time.Time `json:"publishedAt"`
	AccountID       uint       `json:"-"`
	Account         *Account   `json:"-"`
	Address1        *string    `json:"address1"`
	Address2        *string    `json:"address2"`
	Address3        *string    `json:"address3"`
	PostalCode      *string    `json:"postalCode"`
	CityID          *uint      `json:"cityID"`
	City            *City      `json:"-"`
	CountryID       *uint      `json:"countryID"`
	Country         *Country   `json:"-"`
	Latitude        float32    `json:"latitude"`
	Longitude       float32    `json:"longitude"`
	Images          []*Image   `json:"images"`
	MinimumStay     *int       `json:"minimumStay"`
	Deposit         *string    `json:"deposit"`
	Cleaning        *string    `json:"cleaning"`
	NearestAirport  *string    `json:"nearestAirport"`
	NearestSubway   *string    `json:"nearestSubway"`
	NearbyLocations *string    `json:"nearbyLocations"`
	Overview        *string    `json:"overview"`
	Amenities       []*Amenity `json:"amenities"`
}
