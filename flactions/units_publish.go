package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// UnitsPublish marks a property as published, allowing it to appear in search results
type UnitsPublish struct{}

// Perform executes the action
func (a *UnitsPublish) Perform(c Context) {
	user := c.CurrentUser()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Diagnostics().AddError(err)
		c.RespondNotFound()
		return
	}

	properties := c.Repository().UserProperties(user).Select("id")

	unit := &flentities.Unit{}
	err = c.Repository().Preload("Amenities").Preload("Images", flentities.Image{Uploaded: true}).
		Where(map[string]interface{}{"id": id}).Where("property_id IN (?)", properties.QueryExpr()).Find(unit).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.RespondNotFound()
		} else {
			c.ServerError(err)
		}
		return
	}

	errs := unit.CanBePublished()
	if len(errs) > 0 {
		c.RespondError(http.StatusUnprocessableEntity, errs...)
		return
	}

	if err = c.Repository().Model(unit).Updates(map[string]interface{}{"PublishedAt": c.Now()}).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, nil)
}
