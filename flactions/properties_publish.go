package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesPublish marks a property as published, allowing it to appear in search results
type PropertiesPublish struct{}

// Perform executes the action
func (a *PropertiesPublish) Perform(c Context) {
	user := c.CurrentUser()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Diagnostics().AddError(err)
		c.RespondNotFound()
		return
	}

	conditions := map[string]interface{}{
		"id":         id,
		"account_id": *user.AccountID,
	}

	property := &flentities.Property{}
	err = c.Repository().Preload("Amenities").Preload("Units").Preload("Images", flentities.Image{Uploaded: true}).
		Find(property, conditions).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.RespondNotFound()
		} else {
			c.ServerError(err)
		}
		return
	}

	errs := property.CanBePublished()
	if len(errs) > 0 {
		c.RespondError(http.StatusUnprocessableEntity, errs...)
		return
	}

	if err = c.Repository().Model(property).Updates(map[string]interface{}{"PublishedAt": c.Now()}).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, nil)
}
