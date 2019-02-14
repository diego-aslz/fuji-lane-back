package flactions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/nerde/fuji-lane-back/flconfig"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flservices"
	"github.com/nerde/fuji-lane-back/fujilane"
)

// ImagesCreateBody is the request body for creating a property image
type ImagesCreateBody struct {
	PropertyID uint   `json:"propertyID"`
	UnitID     uint   `json:"unitID"`
	Name       string `json:"name"`
	Size       int    `json:"size"`
	Type       string `json:"type"`
	Position   int    `json:"position"`
}

// Validate the request body
func (b ImagesCreateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("name", b.Name).Required(),
		flentities.ValidateField("size", b.Size).Min(1).Max(flconfig.Config.MaxImageSizeMB*1024*1024),
		flentities.ValidateField("type", b.Type).Required().Image(),
	)
}

// ImagesCreate returns a pre-signed URL for clients to upload images directly to S3
type ImagesCreate struct {
	flservices.S3Service `json:"-"`
	ImagesCreateBody
	Context
}

// Perform executes the action
func (a *ImagesCreate) Perform() {
	account := a.CurrentAccount()

	collection := "properties"
	id := a.PropertyID
	if a.UnitID > 0 {
		collection = "units"
		id = a.UnitID
	}

	key := fmt.Sprintf("public/%s/%d/images/%s", collection, id, fujilane.GenerateRandomString(30, a.RandomSource()))
	url, err := a.GenerateURLToUploadPublicFile(key, a.Type, a.Size)

	if err != nil {
		a.ServerError(err)
		return
	}

	a.Name = strings.Replace(a.Name, "/", "", -1)
	image := &flentities.Image{
		Name:       a.Name,
		URL:        strings.Split(url, "?")[0],
		Type:       a.Type,
		Size:       a.Size,
		StorageKey: &key,
		Position:   a.Position,
	}

	if a.PropertyID > 0 {
		property := &flentities.Property{}
		err := a.Repository().Find(property, map[string]interface{}{"id": a.PropertyID, "account_id": account.ID}).Error
		if gorm.IsRecordNotFoundError(err) {
			a.RespondError(http.StatusUnprocessableEntity, errors.New("Could not find Property"))
			return
		}
		if err != nil {
			a.ServerError(err)
			return
		}
		image.PropertyID = &property.ID

	} else if a.UnitID > 0 {
		unit := &flentities.Unit{}
		err := a.Repository().Preload("Property").Find(unit, map[string]interface{}{"id": a.UnitID}).Error
		if gorm.IsRecordNotFoundError(err) || unit.Property == nil || unit.Property.AccountID != a.CurrentAccount().ID {
			a.RespondError(http.StatusUnprocessableEntity, errors.New("Could not find Unit"))
			return
		}
		if err != nil {
			a.ServerError(err)
			return
		}
		image.UnitID = &unit.ID

	} else {
		a.RespondError(http.StatusUnprocessableEntity, errors.New("Please provide either a Property or a Unit"))
		return
	}

	if err = a.Repository().Save(image).Error; err != nil {
		a.ServerError(err)
		return
	}

	image.URL = url

	a.Respond(http.StatusOK, image)
}

// NewImagesCreate returns a new instance of the PropertiesImagesNew action
func NewImagesCreate(s3 flservices.S3Service, c Context) *ImagesCreate {
	return &ImagesCreate{S3Service: s3, Context: c}
}
