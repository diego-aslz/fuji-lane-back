package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flviews"
)

// ListingsShow exposes details for a property
type ListingsShow struct{}

// Perform executes the action
func (a *ListingsShow) Perform(c Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Diagnostics().AddError(err)
		c.RespondNotFound()
		return
	}

	publishedNull := map[string]interface{}{"published_at": nil}

	property := &flentities.Property{}

	err = c.Repository().Preload("Amenities").Preload("Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).
		Preload("Units", func(db *gorm.DB) *gorm.DB { return db.Not(publishedNull).Order("units.base_price_cents") }).
		Preload("Units.Images", flentities.Image{Uploaded: true}, imagesDefaultOrder).Preload("Units.Amenities").
		Where(id).Not(publishedNull).Find(property).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.RespondNotFound()
		} else {
			c.ServerError(err)
		}
		return
	}

	similarListings := []*flentities.Property{}
	err = c.Repository().Preload("Images", flentities.Image{Uploaded: true}, imagesDefaultOrder,
		func(db *gorm.DB) *gorm.DB { return db.Limit(1) }).
		Preload("Units", func(db *gorm.DB) *gorm.DB {
			return db.Not(publishedNull).Order("units.base_price_cents").Limit(1)
		}).
		Where(map[string]interface{}{"city_id": property.CityID}).
		Where("id in (?)", c.Repository().Table("units").Select("property_id").Not(publishedNull).QueryExpr()).
		Not(publishedNull).Not(property.ID).Limit(3).Find(&similarListings).Error

	if err != nil {
		c.Diagnostics().Add("similar_listings_error", err.Error())
	}

	c.Respond(http.StatusOK, flviews.NewListing(property, similarListings))
}
