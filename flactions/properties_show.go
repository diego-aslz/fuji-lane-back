package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesShow exposes details for a property
type PropertiesShow struct {
	Context
}

// Perform executes the action
func (a *PropertiesShow) Perform() {
	user := a.CurrentUser()

	id, err := strconv.Atoi(a.Param("id"))
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondNotFound()
		return
	}

	conditions := map[string]interface{}{
		"id":         id,
		"account_id": *user.AccountID,
	}

	property := &flentities.Property{}
	err = a.Repository().Preload("Amenities").Preload("Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).
		Preload("Units.Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).Preload("Units.Amenities").
		Find(property, conditions).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			a.RespondNotFound()
		} else {
			a.ServerError(err)
		}
		return
	}

	a.Respond(http.StatusOK, property)
}

// NewPropertiesShow returns a new PropertiesShow action
func NewPropertiesShow(c Context) Action {
	return &PropertiesShow{c}
}
