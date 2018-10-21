package flactions

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// PropertiesUpdateBody is the request body for creating a property image
type PropertiesUpdateBody struct {
	Name            *string `json:"name"`
	Address1        *string `json:"address1"`
	Address2        *string `json:"address2"`
	Address3        *string `json:"address3"`
	CityID          *int    `json:"cityID"`
	PostalCode      *string `json:"postalCode"`
	MinimumStay     *string `json:"minimumStay"`
	Deposit         *string `json:"deposit"`
	Cleaning        *string `json:"cleaning"`
	NearestAirport  *string `json:"nearestAirport"`
	NearestSubway   *string `json:"nearestSubway"`
	NearbyLocations *string `json:"nearbyLocations"`
	Overview        *string `json:"overview"`
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
	err := c.Repository().Find(property, map[string]interface{}{"id": id, "account_id": account.ID}).Error
	if gorm.IsRecordNotFoundError(err) {
		c.RespondNotFound()
		return
	}
	if err != nil {
		c.ServerError(err)
		return
	}

	updates := map[string]interface{}{}

	fields := map[string]interface{}{
		"name":             a.Name,
		"address1":         a.Address1,
		"address2":         a.Address2,
		"address3":         a.Address3,
		"postal_code":      a.PostalCode,
		"minimum_stay":     a.MinimumStay,
		"deposit":          a.Deposit,
		"cleaning":         a.Cleaning,
		"nearest_airport":  a.NearestAirport,
		"nearest_subway":   a.NearestSubway,
		"nearby_locations": a.NearbyLocations,
		"overview":         a.Overview,
	}
	for field, value := range fields {
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

	if err := c.Repository().Model(property).Updates(updates).Error; err != nil {
		c.ServerError(err)
		return
	}

	c.Respond(http.StatusOK, nil)
}
