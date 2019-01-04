package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// CountriesList lists the available countries
type CountriesList struct {
	Context
}

// Perform executes the action
func (a *CountriesList) Perform() {
	countries := []*flentities.Country{}
	if err := a.Repository().Order("name").Find(&countries).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, countries)
}

// NewCountriesList returns a new CountriesList action
func NewCountriesList(c Context) Action {
	return &CountriesList{c}
}
