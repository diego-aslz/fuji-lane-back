package flactions

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flservices"
)

// PropertiesImagesNew returns a pre-signed URL for clients to upload images directly to S3
type PropertiesImagesNew struct {
	*flservices.S3
}

// Perform executes the action
func (a *PropertiesImagesNew) Perform(c Context) {
	fileName := c.Query("name")

	c.Diagnostics().AddQuoted("file_name", fileName)

	if fileName == "" {
		c.RespondError(http.StatusUnprocessableEntity, errors.New("Please provide a filename"))
		return
	}

	account := c.CurrentAccount()
	if account == nil {
		c.RespondError(http.StatusPreconditionRequired, errors.New("You do not have an owner account"))
		return
	}

	id := c.Param("id")
	property := &flentities.Property{}
	err := c.Repository().Find(property, map[string]interface{}{"id": id, "account_id": account.ID}).Error
	if gorm.IsRecordNotFoundError(err) {
		c.RespondError(http.StatusNotFound, errors.New("Not found"))
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	url, err := a.GenerateURLToUploadPublicFile("properties/" + id + "/images/" + strings.Replace(fileName, "/", "", -1))

	if err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, map[string]string{"url": url})
}

// NewPropertiesImagesNew returns a new instance of the PropertiesImagesNew action
func NewPropertiesImagesNew() *PropertiesImagesNew {
	return &PropertiesImagesNew{flservices.NewS3()}
}
