package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// CountriesList lists the available countries
type CountriesList struct{}

// Perform executes the action
func (a *CountriesList) Perform(c Context) {
	countries := []*flentities.Country{}
	if err := c.Repository().Order("name").Find(&countries).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, countries)
}
