package flentities

import "time"

// Price represents a unit price under specific conditions
type Price struct {
	ID        uint      `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	UnitID    uint      `json:"-"`
	Unit      *Unit     `json:"-"`
	MinNights int       `json:"minNights"`
	Cents     int       `json:"cents"`
}
