package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// CitiesList lists the available countries
type CitiesList struct {
	Context
}

// Perform executes the action
func (a *CitiesList) Perform() {
	cities := []*flentities.City{}
	if err := a.Repository().Order("name").Find(&cities).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, cities)
}

// NewCitiesList returns a new CitiesList action
func NewCitiesList(c Context) Action {
	return &CitiesList{c}
}
