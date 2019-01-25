package flentities

import (
	"errors"
	"time"

	"github.com/nerde/fuji-lane-back/flutils"
)

// Unit represents a type of rentable unit inside a Property. An individual property can have multiple units at
// different prices for rent
type Unit struct {
	ID               uint       `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time  `json:"-"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `json:"-"`
	PublishedAt      *time.Time `json:"publishedAt"`
	EverPublished    bool       `json:"everPublished"`
	PropertyID       uint       `json:"propertyID"`
	Property         *Property  `json:"-"`
	Name             string     `json:"name"`
	Slug             string     `json:"slug"`
	Overview         *string    `json:"overview"`
	Bedrooms         int        `json:"bedrooms"`
	Bathrooms        int        `json:"bathrooms"`
	SizeM2           int        `json:"sizeM2"`
	MaxOccupancy     *int       `json:"maxOccupancy"`
	Count            int        `json:"count"`
	FloorPlanImageID *uint      `json:"-"`
	FloorPlanImage   *Image     `json:"floorPlanImage"`
	Amenities        []*Amenity `json:"amenities"`
	Images           []*Image   `json:"images"`
	Prices           []*Price   `json:"prices"`
}

// BeforeSave to update the slug
func (u *Unit) BeforeSave() error {
	if u.Name != "" {
		u.Slug = generateSlug(u.Name)
	}

	return nil
}

// CanBePublished checks if this unit can be marked as published and start showing up in search results
func (u *Unit) CanBePublished() []error {
	errs := []error{}

	if flutils.IsBlankStr(&u.Name) {
		errs = append(errs, errors.New("Name is required"))
	}

	if u.Bedrooms == 0 {
		errs = append(errs, errors.New("Bedrooms is required"))
	}

	if u.SizeM2 == 0 {
		errs = append(errs, errors.New("Size is required"))
	}

	if u.Count == 0 {
		errs = append(errs, errors.New("Number of Unit Type is required"))
	}

	if u.Amenities == nil || len(u.Amenities) == 0 {
		errs = append(errs, errors.New("At least one amenity is required"))
	}

	if !u.hasUploadedImages() {
		errs = append(errs, errors.New("At least one image is required"))
	}

	if u.basePriceCents() == 0 {
		errs = append(errs, errors.New("Please provide a base unit price"))
	}

	return errs
}

func (u *Unit) basePriceCents() int {
	if u.Prices == nil || len(u.Prices) == 0 {
		return 0
	}

	for _, p := range u.Prices {
		if p.MinNights == 1 {
			return p.Cents
		}
	}

	return 0
}

func (u *Unit) hasUploadedImages() bool {
	for _, i := range u.Images {
		if i.Uploaded {
			return true
		}
	}

	return false
}

// PriceFor returns the price that should be used when booking the given number nights
func (u *Unit) PriceFor(nights int) Price {
	lastMinNights := 0
	price := Price{}

	if u.Prices != nil {
		for _, p := range u.Prices {
			if p.MinNights > lastMinNights && p.MinNights <= nights {
				lastMinNights = p.MinNights
				price = *p
			}
		}
	}

	return price
}
