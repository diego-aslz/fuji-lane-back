package flactions

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// ImagesUploaded marks an image as uploaded
type ImagesUploaded struct {
	Context
}

// Perform the action
func (a *ImagesUploaded) Perform() {
	account := a.CurrentAccount()

	id := a.Param("id")
	image := &flentities.Image{}
	err := a.Repository().Preload("Property").Preload("Unit.Property").Find(image,
		map[string]interface{}{"id": id, "uploaded": false}).Error
	if gorm.IsRecordNotFoundError(err) {
		a.Diagnostics().AddQuoted("reason", "Could not find Image")
		a.RespondNotFound()
		return
	}
	if err != nil {
		a.ServerError(err)
		return
	}

	var imageAccountID uint
	if image.Property != nil {
		imageAccountID = image.Property.AccountID
	} else {
		imageAccountID = image.Unit.Property.AccountID
	}

	if imageAccountID != account.ID {
		a.Diagnostics().AddQuoted("reason", "Image belongs to another account")
		a.RespondNotFound()
		return
	}

	if err = a.Repository().Model(image).Updates(map[string]interface{}{"uploaded": true}).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, image)
}

// NewImagesUploaded returns a new ImagesUploaded action
func NewImagesUploaded(c Context) Action {
	return &ImagesUploaded{c}
}
