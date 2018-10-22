package flactions

import (
	"net/http"
	"strings"

	"github.com/nerde/fuji-lane-back/flconfig"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flservices"
	"github.com/nerde/fuji-lane-back/flutils"
)

// PropertiesImagesCreateBody is the request body for creating a property image
type PropertiesImagesCreateBody struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	Type string `json:"type"`
}

// Validate the request body
func (b PropertiesImagesCreateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("name", b.Name).Required(),
		flentities.ValidateField("size", b.Size).Min(1).Max(flconfig.Config.MaxImageSizeMB*1024*1024),
		flentities.ValidateField("type", b.Type).Required().Image(),
	)
}

// PropertiesImagesCreate returns a pre-signed URL for clients to upload images directly to S3
type PropertiesImagesCreate struct {
	flservices.S3Service
	PropertiesImagesCreateBody
}

// Perform executes the action
func (a *PropertiesImagesCreate) Perform(c Context) {
	account := c.CurrentAccount()

	id := c.Param("id")
	property := &flentities.Property{}
	err := c.Repository().Find(property, map[string]interface{}{"id": id, "account_id": account.ID}).Error
	if gorm.IsRecordNotFoundError(err) {
		c.RespondNotFound()
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	key := flutils.GenerateRandomString(30, c.RandomSource())
	url, err := a.GenerateURLToUploadPublicFile("properties/"+id+"/images/"+key, a.Type, a.Size)

	if err != nil {
		c.ServerError(err)
		return
	}

	a.Name = strings.Replace(a.Name, "/", "", -1)
	image := &flentities.Image{
		Name:       a.Name,
		URL:        strings.Split(url, "?")[0],
		Type:       a.Type,
		Size:       a.Size,
		PropertyID: &property.ID,
	}

	if err = c.Repository().Save(image).Error; err != nil {
		c.ServerError(err)
		return
	}

	image.URL = url

	c.Respond(http.StatusOK, image)
}

// NewPropertiesImagesCreate returns a new instance of the PropertiesImagesNew action
func NewPropertiesImagesCreate(s3 flservices.S3Service) *PropertiesImagesCreate {
	return &PropertiesImagesCreate{S3Service: s3}
}
