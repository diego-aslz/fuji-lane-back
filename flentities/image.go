package flentities

import (
	"time"
)

// Image we support
type Image struct {
	ID         uint      `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Size       int       `json:"size"`
	URL        string    `json:"url"`
	Uploaded   bool      `json:"uploaded"`
	PropertyID int       `json:"propertyID"`
	Property   Property  `json:"-"`
	Unit       Unit      `json:"-"`
}
