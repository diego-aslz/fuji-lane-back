package flactions

import (
	"net/http"

	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesList lists user properties
type PropertiesList struct {
	Context
}

// Perform executes the action
func (a *PropertiesList) Perform() {
	user := a.CurrentUser()

	properties := []*flentities.Property{}
	err := a.Repository().Order("name").Preload("Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).
		Preload("Units").Find(&properties, map[string]interface{}{"account_id": user.AccountID}).Error
	if err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, properties)
}

// NewPropertiesList returns a new PropertiesList action
func NewPropertiesList(c Context) Action {
	return &PropertiesList{c}
}
