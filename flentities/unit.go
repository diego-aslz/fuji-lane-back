package flentities

import "time"

// Unit represents a type of rentable unit inside a Property. An individual property can have multiple units at
// different prices for rent
type Unit struct {
	ID                     uint       `gorm:"primary_key" json:"id"`
	CreatedAt              time.Time  `json:"-"`
	UpdatedAt              time.Time  `json:"-"`
	DeletedAt              *time.Time `sql:"index" json:"-"`
	PropertyID             *uint      `json:"propertyID"`
	Property               *Property  `json:"-"`
	Name                   string     `json:"name"`
	Bedrooms               int        `json:"bedrooms"`
	SizeM2                 *int       `json:"sizeM2"`
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
