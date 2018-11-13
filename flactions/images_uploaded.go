package flactions

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// ImagesUploaded marks an image as uploaded
type ImagesUploaded struct{}

// Perform the action
func (a *ImagesUploaded) Perform(c Context) {
	account := c.CurrentAccount()

	id := c.Param("id")
	image := &flentities.Image{}
	err := c.Repository().Preload("Property").Preload("Unit.Property").Find(image,
		map[string]interface{}{"id": id, "uploaded": false}).Error
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

	if err = c.Repository().Model(image).Updates(map[string]interface{}{"uploaded": true}).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, image)
}
