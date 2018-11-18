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
	PropertyID *uint     `json:"-"`
	Property   *Property `json:"-"`
	UnitID     *uint     `json:"-"`
	Unit       *Unit     `json:"-"`
	Position   int       `json:"position"`
}

// BelongsTo determines if this image belongs to the given account by its ID
func (i Image) BelongsTo(accountID uint) bool {
	return i.Property != nil && i.Property.AccountID == accountID || i.Unit != nil && i.Unit.Property != nil && i.Unit.Property.AccountID == accountID
}
