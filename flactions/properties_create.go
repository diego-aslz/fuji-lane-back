package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesCreate creates properties that can hold units
type PropertiesCreate struct {
	Context
}

// Perform executes the action
func (a *PropertiesCreate) Perform() {
	user := a.CurrentUser()

	property := &flentities.Property{AccountID: *user.AccountID}

	if err := a.Repository().Create(property).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusCreated, property)
}

// NewPropertiesCreate returns a new PropertiesCreate action
func NewPropertiesCreate(c Context) Action {
	return &PropertiesCreate{c}
}
