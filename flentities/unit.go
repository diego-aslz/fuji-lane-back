package flentities

import "github.com/jinzhu/gorm"

// Unit represents a type of rentable unit inside a Property. An individual property can have multiple units at
// different prices for rent
type Unit struct {
	gorm.Model
	PropertyID             *uint     `json:"-"`
	Property               *Property `json:"-"`
	Name                   string    `json:"name"`
	Bedrooms               int       `json:"bedrooms"`
	SizeM2                 *int      `json:"sizeM2"`
	MaxOccupancy           *int      `json:"maxOccupancy"`
	Count                  int       `json:"count"`
	BasePriceCents         *int      `json:"basePriceCents"`
	OneNightPriceCents     *int      `json:"oneNightPriceCents"`
	OneWeekPriceCents      *int      `json:"oneWeekPriceCents"`
	ThreeMonthsPriceCents  *int      `json:"threeMonthsPriceCents"`
	SixMonthsPriceCents    *int      `json:"sixMonthsPriceCents"`
	TwelveMonthsPriceCents *int      `json:"twelveMonthsPriceCents"`
	FloorPlanImageID       *uint     `json:"-"`
	FloorPlanImage         *Image    `json:"floorPlanImage"`
}
