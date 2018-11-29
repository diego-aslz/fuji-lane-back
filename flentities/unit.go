package flentities

import (
	"errors"
	"time"

	"github.com/nerde/fuji-lane-back/flutils"
)

// Unit represents a type of rentable unit inside a Property. An individual property can have multiple units at
// different prices for rent
type Unit struct {
	ID                     uint       `gorm:"primary_key" json:"id"`
	CreatedAt              time.Time  `json:"-"`
	UpdatedAt              time.Time  `json:"-"`
	DeletedAt              *time.Time `json:"-"`
	PublishedAt            *time.Time `json:"publishedAt"`
	PropertyID             *uint      `json:"propertyID"`
	Property               *Property  `json:"-"`
	Name                   string     `json:"name"`
	Bedrooms               int        `json:"bedrooms"`
	SizeM2                 int        `json:"sizeM2"`
	MaxOccupancy           *int       `json:"maxOccupancy"`
	Count                  int        `json:"count"`
	BasePriceCents         *int       `json:"basePriceCents"`
	OneNightPriceCents     *int       `json:"oneNightPriceCents"`
	OneWeekPriceCents      *int       `json:"oneWeekPriceCents"`
	ThreeMonthsPriceCents  *int       `json:"threeMonthsPriceCents"`
	SixMonthsPriceCents    *int       `json:"sixMonthsPriceCents"`
	TwelveMonthsPriceCents *int       `json:"twelveMonthsPriceCents"`
	FloorPlanImageID       *uint      `json:"-"`
	FloorPlanImage         *Image     `json:"floorPlanImage"`
	Amenities              []*Amenity `json:"amenities"`
	Images                 []*Image   `json:"images"`
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

	if len(u.Images) < 1 {
		errs = append(errs, errors.New("At least one image is required"))
	}

	return errs
}
