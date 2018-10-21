package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// AmenityTypesList lists the amenity types for the given target
type AmenityTypesList struct{}

// Perform executes the action
func (a *AmenityTypesList) Perform(c Context) {
	target := c.Param("target")
	c.Diagnostics().Add("target", target)

	if target != "properties" && target != "units" {
		c.RespondNotFound()
		return
	}

	types := []*flentities.AmenityType{}

	for _, amType := range flentities.AmenityTypes {
		if target == "properties" && amType.ForProperties || target == "units" && amType.ForUnits {
			types = append(types, amType)
		}
	}

	c.Respond(http.StatusOK, types)
}
