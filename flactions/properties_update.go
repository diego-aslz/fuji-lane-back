package flactions

import (
	"errors"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/optional"
)

// PropertiesUpdateBody is the request body for updating a property image
type PropertiesUpdateBody struct {
	Name            optional.String  `json:"name"`
	Address1        optional.String  `json:"address1"`
	Address2        optional.String  `json:"address2"`
	Address3        optional.String  `json:"address3"`
	CityID          optional.Uint    `json:"cityID"`
	Latitude        optional.Float32 `json:"latitude"`
	Longitude       optional.Float32 `json:"longitude"`
	PostalCode      optional.String  `json:"postalCode"`
	MinimumStay     optional.Int     `json:"minimumStay"`
	Deposit         optional.String  `json:"deposit"`
	Cleaning        optional.String  `json:"cleaning"`
	NearestAirport  optional.String  `json:"nearestAirport"`
	NearestSubway   optional.String  `json:"nearestSubway"`
	NearbyLocations optional.String  `json:"nearbyLocations"`
	Overview        optional.String  `json:"overview"`
	bodyWithAmenities
}

// Validate the request body
func (b *PropertiesUpdateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("overview", b.Overview.Value).HTML(),
	)
}

// PropertiesUpdate returns a pre-signed URL for clients to upload images directly to S3
type PropertiesUpdate struct {
	PropertiesUpdateBody
	Context
}

// Perform executes the action
func (a *PropertiesUpdate) Perform() {
	account := a.CurrentAccount()

	id := a.Param("id")
	property := &flentities.Property{}
	conditions := map[string]interface{}{"id": id, "account_id": account.ID}
	err := a.Repository().Preload("Amenities").Find(property, conditions).Error
	if gorm.IsRecordNotFoundError(err) {
		a.RespondNotFound()
		return
	}
	if err != nil {
		a.ServerError(err)
		return
	}

	optional.Update(a.PropertiesUpdateBody, property)

	if a.CityID.Set && a.CityID.Value != nil {
		city := &flentities.City{}
		city.ID = uint(*a.CityID.Value)
		err := a.Repository().Find(city).Error

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
	}

	a.Repository().Transaction(func(tx *flentities.Repository) {
		amenitiesToDelete, amenitiesToCreate := a.amenitiesDiff(property.Amenities)

		for _, am := range amenitiesToDelete {
			if err := tx.Delete(am).Error; err != nil {
				tx.Rollback()
				a.ServerError(err)
				return
			}
		}

		for _, am := range amenitiesToCreate {
			am.PropertyID = &property.ID
			if err := tx.Create(am).Error; err != nil {
				tx.Rollback()
				a.ServerError(err)
				return
			}
		}

		if err := tx.Save(property).Error; err != nil {
			tx.Rollback()

			if flentities.IsUniqueConstraintViolation(err) &&
				(strings.Index(err.Error(), "_name") > -1 || strings.Index(err.Error(), "_slug") > -1) {
				a.RespondError(http.StatusUnprocessableEntity, errors.New("Name is already in use"))
				return
			}

			a.ServerError(err)
			return
		}

		if err := tx.Commit().Error; err != nil {
			a.ServerError(err)
			return
		}

		a.Respond(http.StatusOK, property)
	})
}

// NewPropertiesUpdate returns a new PropertiesUpdate action
func NewPropertiesUpdate(c Context) Action {
	return &PropertiesUpdate{Context: c}
}
