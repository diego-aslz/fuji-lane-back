package flactions

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/optional"
)

// PropertiesCreateBody is the request body for creating a property image
type PropertiesCreateBody struct {
	Name       string           `json:"name"`
	Overview   optional.String  `json:"overview"`
	Address1   string           `json:"address1"`
	Address2   optional.String  `json:"address2"`
	Address3   optional.String  `json:"address3"`
	CityID     uint             `json:"cityID"`
	PostalCode optional.String  `json:"postalCode"`
	Latitude   optional.Float32 `json:"latitude"`
	Longitude  optional.Float32 `json:"longitude"`
}

// Validate the request body
func (b *PropertiesCreateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("name", b.Name).Required(),
		flentities.ValidateField("overview", b.Overview.Value).HTML(),
		flentities.ValidateField("address", b.Address1).Required(),
		flentities.ValidateField("city", b.CityID).Required(),
	)
}

// PropertiesCreate creates properties that can hold units
type PropertiesCreate struct {
	Context
	PropertiesCreateBody
}

// Perform executes the action
func (a *PropertiesCreate) Perform() {
	user := a.CurrentUser()

	property := &flentities.Property{
		AccountID: *user.AccountID,
		Name:      a.Name,
		Address1:  a.Address1,
		CityID:    a.CityID,
	}

	city := &flentities.City{}
	err := a.Repository().Find(city, a.CityID).Error

	if gorm.IsRecordNotFoundError(err) {
		a.RespondError(http.StatusUnprocessableEntity, errors.New("Invalid City"))
		return
	}
	if err != nil {
		a.ServerError(err)
		return
	}

	property.CityID = city.ID
	property.CountryID = city.CountryID

	optional.Update(a.PropertiesCreateBody, property)

	if err := a.Repository().Create(property).Error; err != nil {
		a.ServerError(err)
		return
	}

	a.Respond(http.StatusCreated, property)
}

// NewPropertiesCreate returns a new PropertiesCreate action
func NewPropertiesCreate(c Context) Action {
	return &PropertiesCreate{Context: c}
}
