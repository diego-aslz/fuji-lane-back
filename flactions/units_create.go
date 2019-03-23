package flactions

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/optional"
)

// UnitsCreateBody is the representation of the payload for creating a Unit
type UnitsCreateBody struct {
	PropertyID   uint            `json:"propertyID"`
	Name         string          `json:"name"`
	Overview     optional.String `json:"overview"`
	Bedrooms     int             `json:"bedrooms"`
	SizeM2       int             `json:"sizeM2"`
	SizeFT2      int             `json:"sizeFT2"`
	MaxOccupancy optional.Int    `json:"maxOccupancy"`
	Count        int             `json:"count"`
}

// Validate the request body
func (b *UnitsCreateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("Property", b.PropertyID).Required(),
		flentities.ValidateField("Name", b.Name).Required(),
		flentities.ValidateField("Bedrooms", b.Bedrooms).Required(),
		flentities.ValidateField("Size in m²", b.SizeM2).Required(),
		flentities.ValidateField("Size in ft²", b.SizeFT2).Required(),
		flentities.ValidateField("Number of unit type", b.Count).Required(),
		flentities.ValidateField("Overview", b.Overview.Value).HTML(),
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

	conditions := map[string]interface{}{"id": a.PropertyID, "account_id": a.CurrentUser().AccountID}
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
		Property: property,
		Name:     a.Name,
		Bedrooms: a.Bedrooms,
		SizeM2:   a.SizeM2,
		SizeFT2:  a.SizeFT2,
		Count:    a.Count,
	}

	optional.Update(a.UnitsCreateBody, unit)

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
