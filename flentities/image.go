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
	URL        string    `json:"url"`
	Uploaded   bool      `json:"uploaded"`
	PropertyID int       `json:"propertyID"`
	Property   Property  `json:"-"`
}
