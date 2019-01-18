package flactions

import (
	"errors"
	"net/http"
	"strings"

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
		flentities.ValidateField("overview", b.Overview).HTML(),
	)
}

// UnitsCreate creates a new Unit
type UnitsCreate struct {
	UnitsCreateBody
	Context
}

// Perform executes the action
func (a *UnitsCreate) Perform() {
	property := &flentities.Property{}

	conditions := map[string]interface{}{"id": a.PropertyID, "account_id": a.CurrentAccount().ID}
	err := a.Repository().Find(property, conditions).Error
	if gorm.IsRecordNotFoundError(err) {
		a.Diagnostics().AddQuoted("reason", "Could not find property")
		a.RespondNotFound()
		return
	}
	if err != nil {
		a.ServerError(err)
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

	if err := a.Repository().Create(unit).Error; err != nil {
		if flentities.IsUniqueConstraintViolation(err) && strings.Index(err.Error(), "_slug") > -1 {
			a.RespondError(http.StatusUnprocessableEntity, errors.New("Name is already in use"))
			return
		}

		a.ServerError(err)
		return
	}

	a.Respond(http.StatusCreated, unit)
}

// NewUnitsCreate returns a new UnitsCreate action
func NewUnitsCreate(c Context) Action {
	return &UnitsCreate{Context: c}
}
