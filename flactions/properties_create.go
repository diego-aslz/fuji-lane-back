package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesCreate creates properties that can hold units
type PropertiesCreate struct{}

// Perform executes the action
func (a *PropertiesCreate) Perform(c Context) {
	user := c.CurrentUser()

	property := &flentities.Property{AccountID: *user.AccountID}

	if err := c.Repository().Create(property).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusCreated, property)
}
