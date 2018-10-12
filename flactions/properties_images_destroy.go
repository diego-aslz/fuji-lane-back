package flactions

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flservices"
)

// PropertiesImagesDestroy marks an image as uploaded
type PropertiesImagesDestroy struct {
	flservices.S3Service
}

// Perform the action
func (a *PropertiesImagesDestroy) Perform(c Context) {
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
	err = c.Repository().Find(image, map[string]interface{}{"id": id, "property_id": property.ID}).Error
	if gorm.IsRecordNotFoundError(err) {
		c.Diagnostics().AddQuoted("reason", "Could not find Image")
		c.RespondNotFound()
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	if err = a.DeleteFile(image.URL); err != nil {
		c.ServerError(err)
		return
	}

	if err = c.Repository().Delete(image).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, map[string]interface{}{})
}

// NewPropertiesImagesDestroy returns a new instance of the PropertiesImagesDestroy action
func NewPropertiesImagesDestroy(s3 flservices.S3Service) *PropertiesImagesDestroy {
	return &PropertiesImagesDestroy{s3}
}
