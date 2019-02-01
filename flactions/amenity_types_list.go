package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// AmenityTypesList lists the amenity types for the given target
type AmenityTypesList struct {
	Context
}

// Perform executes the action
func (a *AmenityTypesList) Perform() {
	target := a.Param("target")
	a.Diagnostics().Add("target", target)

	var types []*flentities.AmenityType

	if target == "properties" {
		types = flentities.PropertyAmenityTypes
	} else if target == "units" {
		types = flentities.UnitAmenityTypes
	} else {
		a.RespondNotFound()
		return
	}

	a.Respond(http.StatusOK, types)
}

// NewAmenityTypesList returns a new AmenityTypesList action
func NewAmenityTypesList(c Context) Action {
	return &AmenityTypesList{c}
}
