package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesPublish marks a property as published, allowing it to appear in search results
type PropertiesPublish struct {
	Context
}

// Perform executes the action
func (a *PropertiesPublish) Perform() {
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
	err = a.Repository().Preload("Amenities").Preload("Units").Preload("Images", flentities.Image{Uploaded: true}).
		Find(property, conditions).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			a.RespondNotFound()
		} else {
			a.ServerError(err)
		}
		return
	}

	errs := property.CanBePublished()
	if len(errs) > 0 {
		a.RespondError(http.StatusUnprocessableEntity, errs...)
		return
	}

	updates := map[string]interface{}{"PublishedAt": a.Now(), "EverPublished": true}
	if err = a.Repository().Model(property).Updates(updates).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, nil)
}

// NewPropertiesPublish returns a new PropertiesPublish action
func NewPropertiesPublish(c Context) Action {
	return &PropertiesPublish{c}
}
