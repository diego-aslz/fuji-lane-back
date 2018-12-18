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
}

// Perform executes the action
func (a *PropertiesUpdate) Perform(c Context) {
	account := c.CurrentAccount()

	id := c.Param("id")
	property := &flentities.Property{}
	conditions := map[string]interface{}{"id": id, "account_id": account.ID}
	err := c.Repository().Preload("Amenities").Find(property, conditions).Error
	if gorm.IsRecordNotFoundError(err) {
		c.RespondNotFound()
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	updates := map[string]interface{}{}
	for field, value := range a.bodyMap() {
		if value != nil {
			updates[field] = value
		}
	}

	if a.CityID != nil {
		city := &flentities.City{}
		city.ID = uint(*a.CityID)
		err := c.Repository().Find(city).Error

		if gorm.IsRecordNotFoundError(err) {
			c.RespondError(http.StatusUnprocessableEntity, errors.New("Invalid City"))
			return
		}
		if err != nil {
			c.ServerError(err)
			return
		}

		updates["city_id"] = city.ID
		updates["country_id"] = city.CountryID
	}

	if a.Latitude != nil {
		updates["latitude"] = *a.Latitude
	}

	if a.Longitude != nil {
		updates["longitude"] = *a.Longitude
	}

	if a.MinimumStay != nil {
		updates["MinimumStay"] = *a.MinimumStay
	}

	c.Repository().Transaction(func(tx *flentities.Repository) {
		amenitiesToDelete, amenitiesToCreate := a.amenitiesDiff(property.Amenities)

		for _, am := range amenitiesToDelete {
			if err := tx.Delete(am).Error; err != nil {
				tx.Rollback()
				c.ServerError(err)
				return
			}
		}

		for _, am := range amenitiesToCreate {
			am.PropertyID = &property.ID
			if err := tx.Create(am).Error; err != nil {
				tx.Rollback()
				c.ServerError(err)
				return
			}
		}

		if err := tx.Model(property).Updates(updates).Error; err != nil {
			tx.Rollback()
			c.ServerError(err)
			return
		}

		if err := tx.Commit().Error; err != nil {
			c.ServerError(err)
			return
		}

		c.Respond(http.StatusOK, property)
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
