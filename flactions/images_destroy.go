package flactions

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flservices"
)

// ImagesDestroy destroys an image
type ImagesDestroy struct {
	flservices.S3Service
}

// Perform the action
func (a *ImagesDestroy) Perform(c Context) {
	account := c.CurrentAccount()

	id := c.Param("id")
	image := &flentities.Image{}
	err := c.Repository().Preload("Property").Preload("Unit.Property").Find(image, map[string]interface{}{"id": id}).Error
	if gorm.IsRecordNotFoundError(err) {
		c.Diagnostics().AddQuoted("reason", "Could not find Image")
		c.RespondNotFound()
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	var imageAccountID uint
	if image.Property != nil {
		imageAccountID = image.Property.AccountID
	} else {
		imageAccountID = image.Unit.Property.AccountID
	}

	if imageAccountID != account.ID {
		c.Diagnostics().AddQuoted("reason", "Image belongs to another account")
		c.RespondNotFound()
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

// NewImagesDestroy returns a new instance of the PropertiesImagesDestroy action
func NewImagesDestroy(s3 flservices.S3Service) *ImagesDestroy {
	return &ImagesDestroy{s3}
}
