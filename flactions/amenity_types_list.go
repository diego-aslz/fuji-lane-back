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

	if target != "properties" && target != "units" {
		a.RespondNotFound()
		return
	}

	types := []*flentities.AmenityType{}

	for _, amType := range flentities.AmenityTypes {
		if target == "properties" && amType.ForProperties || target == "units" && amType.ForUnits {
			types = append(types, amType)
		}
	}

	a.Respond(http.StatusOK, types)
}

// NewAmenityTypesList returns a new AmenityTypesList action
func NewAmenityTypesList(c Context) Action {
	return &AmenityTypesList{c}
}
