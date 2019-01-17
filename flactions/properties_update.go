package flactions

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesUpdateBody is the request body for creating a property image
type PropertiesUpdateBody struct {
	Name            *string  `json:"name"`
	Address1        *string  `json:"address1"`
	Address2        *string  `json:"address2"`
	Address3        *string  `json:"address3"`
	CityID          *uint    `json:"cityID"`
	Latitude        *float32 `json:"latitude"`
	Longitude       *float32 `json:"longitude"`
	PostalCode      *string  `json:"postalCode"`
	MinimumStay     *int     `json:"minimumStay"`
	Deposit         *string  `json:"deposit"`
	Cleaning        *string  `json:"cleaning"`
	NearestAirport  *string  `json:"nearestAirport"`
	NearestSubway   *string  `json:"nearestSubway"`
	NearbyLocations *string  `json:"nearbyLocations"`
	Overview        *string  `json:"overview"`
	bodyWithAmenities
}

// Validate the request body
func (b *PropertiesUpdateBody) Validate() []error {
	return flentities.ValidateFields(
		flentities.ValidateField("overview", b.Overview).HTML(),
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

	if a.Name != nil {
		property.Name = a.Name
	}

	if a.Address1 != nil {
		property.Address1 = a.Address1
	}

	if a.Address2 != nil {
		property.Address2 = a.Address2
	}

	if a.Address3 != nil {
		property.Address3 = a.Address3
	}

	if a.PostalCode != nil {
		property.PostalCode = a.PostalCode
	}

	if a.Latitude != nil {
		property.Latitude = *a.Latitude
	}

	if a.Longitude != nil {
		property.Longitude = *a.Longitude
	}

	if a.MinimumStay != nil {
		property.MinimumStay = a.MinimumStay
	}

	if a.Deposit != nil {
		property.Deposit = a.Deposit
	}

	if a.Cleaning != nil {
		property.Cleaning = a.Cleaning
	}

	if a.NearestAirport != nil {
		property.NearestAirport = a.NearestAirport
	}

	if a.NearestSubway != nil {
		property.NearestSubway = a.NearestSubway
	}

	if a.NearbyLocations != nil {
		property.NearbyLocations = a.NearbyLocations
	}

	if a.Overview != nil {
		property.Overview = a.Overview
	}

	if a.CityID != nil {
		city := &flentities.City{}
		city.ID = uint(*a.CityID)
		err := a.Repository().Find(city).Error

		if gorm.IsRecordNotFoundError(err) {
			a.RespondError(http.StatusUnprocessableEntity, errors.New("Invalid City"))
			return
		}
		if err != nil {
			a.ServerError(err)
			return
		}

		property.CityID = &city.ID
		property.CountryID = &city.CountryID
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

func (a *PropertiesUpdate) bodyMap() map[string]*string {
	return map[string]*string{
		"name":             a.Name,
		"address1":         a.Address1,
		"address2":         a.Address2,
		"address3":         a.Address3,
		"postal_code":      a.PostalCode,
		"deposit":          a.Deposit,
		"cleaning":         a.Cleaning,
		"nearest_airport":  a.NearestAirport,
		"nearest_subway":   a.NearestSubway,
		"nearby_locations": a.NearbyLocations,
		"overview":         a.Overview,
	}
}

// NewPropertiesUpdate returns a new PropertiesUpdate action
func NewPropertiesUpdate(c Context) Action {
	return &PropertiesUpdate{Context: c}
}
