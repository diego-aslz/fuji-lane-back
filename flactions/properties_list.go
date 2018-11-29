package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesList lists user properties
type PropertiesList struct{}

// Perform executes the action
func (a *PropertiesList) Perform(c Context) {
	user := c.CurrentUser()

	properties := []*flentities.Property{}
	err := c.Repository().Order("name").Preload("Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).
		Preload("Units").Find(&properties, map[string]interface{}{"account_id": user.AccountID}).Error
	if err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, properties)
}
