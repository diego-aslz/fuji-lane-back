package flactions

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// UnitsCreateBody is the representation of the payload for creating a Unit
type UnitsCreateBody struct {
	PropertyID   uint   `json:"propertyID"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	Bedrooms     int    `json:"bedrooms"`
	SizeM2       int    `json:"sizeM2"`
	MaxOccupancy int    `json:"maxOccupancy"`
	Count        int    `json:"count"`
}

// Validate the request body
func (b *UnitsCreateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("property", b.PropertyID).Required(),
		flentities.ValidateField("name", b.Name).Required(),
		flentities.ValidateField("bedrooms", b.Bedrooms).Required(),
		flentities.ValidateField("size", b.SizeM2).Required(),
		flentities.ValidateField("number of unit type", b.Count).Required(),
	)
}

// UnitsCreate creates a new Unit
type UnitsCreate struct {
	UnitsCreateBody
}

// Perform executes the action
func (a *UnitsCreate) Perform(c Context) {
	property := &flentities.Property{}

	conditions := map[string]interface{}{"id": a.PropertyID, "account_id": c.CurrentAccount().ID}
	err := c.Repository().Find(property, conditions).Error
	if gorm.IsRecordNotFoundError(err) {
		c.Diagnostics().AddQuoted("reason", "Could not find property")
		c.RespondNotFound()
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	unit := &flentities.Unit{
		Property:     property,
		Name:         a.Name,
		Bedrooms:     a.Bedrooms,
		SizeM2:       a.SizeM2,
		MaxOccupancy: &a.MaxOccupancy,
		Count:        a.Count,
	}

	if a.Overview != "" {
		unit.Overview = &a.Overview
	}

	if err := c.Repository().Create(unit).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusCreated, unit)
}
