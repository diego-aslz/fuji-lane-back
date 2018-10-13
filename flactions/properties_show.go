package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesShow creates properties that can hold units
type PropertiesShow struct{}

// Perform executes the action
func (a *PropertiesShow) Perform(c Context) {
	user := c.CurrentUser()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.RespondNotFound()
		return
	}

	conditions := map[string]interface{}{
		"id":         id,
		"account_id": *user.AccountID,
	}

	property := &flentities.Property{}
	err = c.Repository().Preload("Images", flentities.Image{Uploaded: true}).Find(property, conditions).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.RespondNotFound()
		} else {
			c.ServerError(err)
		}
		return
	}

	c.Respond(http.StatusOK, property)
}
