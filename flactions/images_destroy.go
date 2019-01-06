package flactions

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flservices"
)

// ImagesDestroy destroys an image
type ImagesDestroy struct {
	flservices.S3Service `json:"-"`
	Context
}

// Perform the action
func (a *ImagesDestroy) Perform() {
	account := a.CurrentAccount()

	id := a.Param("id")
	image := &flentities.Image{}
	err := a.Repository().Preload("Property").Preload("Unit.Property").Find(image, map[string]interface{}{"id": id}).Error
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

	if err = a.Repository().RemoveImage(image); err != nil {
		a.ServerError(err)
		return
	}

	if err = a.DeleteFile(image.URL); err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusOK, map[string]interface{}{})
}

// NewImagesDestroy returns a new instance of the PropertiesImagesDestroy action
func NewImagesDestroy(s3 flservices.S3Service, c Context) *ImagesDestroy {
	return &ImagesDestroy{S3Service: s3, Context: c}
}
