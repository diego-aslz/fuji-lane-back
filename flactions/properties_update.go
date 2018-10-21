package flactions

import (
	"errors"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
)

// AmenityBody is the payload for an Amenity a Property or Unit can have
type AmenityBody struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	isCreated bool
}

func (ab *AmenityBody) matches(a *flentities.Amenity) bool {
	return ab.matchesDefined(a) || ab.matchesCustom(a)
}

func (ab *AmenityBody) matchesDefined(a *flentities.Amenity) bool {
	return a.Type != "custom" && a.Type == ab.Type
}

func (ab *AmenityBody) matchesCustom(a *flentities.Amenity) bool {
	return a.Type == "custom" && a.Name != nil && ab.Name == *a.Name
}

// PropertiesUpdateBody is the request body for creating a property image
type PropertiesUpdateBody struct {
	Name            *string         `json:"name"`
	Address1        *string         `json:"address1"`
	Address2        *string         `json:"address2"`
	Address3        *string         `json:"address3"`
	CityID          *int            `json:"cityID"`
	PostalCode      *string         `json:"postalCode"`
	MinimumStay     *string         `json:"minimumStay"`
	Deposit         *string         `json:"deposit"`
	Cleaning        *string         `json:"cleaning"`
	NearestAirport  *string         `json:"nearestAirport"`
	NearestSubway   *string         `json:"nearestSubway"`
	NearbyLocations *string         `json:"nearbyLocations"`
	Overview        *string         `json:"overview"`
	Amenities       *[]*AmenityBody `json:"amenities"`
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

	c.Repository().Transaction(func(tx *flentities.Repository) {
		amenitiesToDelete, amenitiesToCreate := a.amenitiesDiff(property)

		for _, am := range amenitiesToDelete {
			if err := tx.Delete(am).Error; err != nil {
				tx.Rollback()
				c.ServerError(err)
				return
			}
		}

		for _, am := range amenitiesToCreate {
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

		c.Respond(http.StatusOK, nil)
	})
}

func (a *PropertiesUpdate) bodyMap() map[string]*string {
	return map[string]*string{
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
}

func (a *PropertiesUpdate) amenitiesDiff(property *flentities.Property) (
	amenitiesToDelete []*flentities.Amenity, amenitiesToCreate []*flentities.Amenity) {

	if a.Amenities == nil {
		return
	}

	// Checking which amenities were removed by the user so we can delete them from the database
	for _, am := range property.Amenities {
		removedByUser := true
		for _, ab := range *a.Amenities {
			if ab.matches(am) {
				removedByUser = false
				ab.isCreated = true
				break
			}
		}

		if removedByUser {
			amenitiesToDelete = append(amenitiesToDelete, am)
		}
	}

	for _, ab := range *a.Amenities {
		// Skipping amenities that are already in the database or invalid
		if ab.isCreated || !flentities.IsValidAmenity(ab.Type, ab.Name) {
			continue
		}

		// Skipping duplicated amenities
		duplicated := false
		for _, am := range amenitiesToCreate {
			if ab.matches(am) {
				duplicated = true
				break
			}
		}
		if duplicated {
			continue
		}

		am := &flentities.Amenity{PropertyID: &property.ID, Type: ab.Type, Name: &ab.Name}
		if ab.Type != "custom" {
			am.Name = nil
		}

		amenitiesToCreate = append(amenitiesToCreate, am)
	}

	return
}
