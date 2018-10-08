package flactions

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesImagesUploaded marks an image as uploaded
type PropertiesImagesUploaded struct{}

// Perform the action
func (a *PropertiesImagesUploaded) Perform(c Context) {
	user := c.CurrentUser()

	id := c.Param("property_id")
	property := &flentities.Property{}
	err := c.Repository().Find(property, map[string]interface{}{"id": id, "account_id": user.AccountID}).Error
	if gorm.IsRecordNotFoundError(err) {
		c.Diagnostics().AddQuoted("reason", "Could not find Property")
		c.RespondNotFound()
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	id = c.Param("id")
	image := &flentities.Image{}
	err = c.Repository().Find(image, map[string]interface{}{"id": id, "property_id": property.ID, "uploaded": false}).Error
	if gorm.IsRecordNotFoundError(err) {
		c.Diagnostics().AddQuoted("reason", "Could not find Image")
		c.RespondNotFound()
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	if err = c.Repository().Model(image).Updates(map[string]interface{}{"uploaded": true}).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, image)
}
