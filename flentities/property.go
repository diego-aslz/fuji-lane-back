package flentities

import (
	"errors"
	"time"

	"github.com/nerde/fuji-lane-back/fujilane"
)

// Property contains address and can have multiple units that can be booked
type Property struct {
	ID              uint       `gorm:"primary_key" json:"id"`
	CreatedAt       time.Time  `json:"-"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"-"`
	Name            *string    `json:"name"`
	Slug            *string    `json:"slug"`
	PublishedAt     *time.Time `json:"publishedAt"`
	EverPublished   bool       `json:"everPublished"`
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
	Units           []*Unit    `json:"units"`
}

// BeforeSave to update the slug
func (p *Property) BeforeSave() error {
	if p.Name != nil && *p.Name != "" {
		slug := generateSlug(*p.Name)
		p.Slug = &slug
	}

	return nil
}

// CanBePublished checks if this property can be marked as published and start showing up in search results
func (p *Property) CanBePublished() []error {
	errs := []error{}

	if fujilane.IsBlankStr(p.Name) {
		errs = append(errs, errors.New("Name is required"))
	}

	if p.isMissingAddress() {
		errs = append(errs, errors.New("Address is incomplete"))
	}

	if p.Amenities == nil || len(p.Amenities) == 0 {
		errs = append(errs, errors.New("At least one amenity is required"))
	}

	if len(p.Images) < 1 {
		errs = append(errs, errors.New("At least one image is required"))
	}

	return errs
}

func (p *Property) isMissingAddress() bool {
	return fujilane.IsBlankStr(p.Address1) || fujilane.IsBlankStr(p.PostalCode) || fujilane.IsBlankUint(p.CityID) ||
		p.Latitude == 0 || p.Longitude == 0
}
