package flactions

import (
	"errors"
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesCreate creates properties that can hold units
type PropertiesCreate struct{}

// Perform executes the action
func (a *PropertiesCreate) Perform(c Context) {
	user := c.CurrentUser()

	if user.AccountID == nil {
		c.RespondError(http.StatusUnprocessableEntity, errors.New("You need a company account"))
		return
	}

	property := &flentities.Property{AccountID: *user.AccountID, StateID: flentities.PropertyStateDraft}

	if err := c.Repository().Create(property).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusCreated, property)
}

// NewPropertiesCreate creates new PropertiesCreate instances
func NewPropertiesCreate() Action {
	return &PropertiesCreate{}
}
