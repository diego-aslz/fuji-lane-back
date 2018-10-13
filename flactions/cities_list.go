package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// CitiesList lists the available countries
type CitiesList struct{}

// Perform executes the action
func (a *CitiesList) Perform(c Context) {
	cities := []*flentities.City{}
	if err := c.Repository().Order("name").Find(&cities).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, cities)
}
