package flactions

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flviews"
)

// ListingsShow exposes details for a property
type ListingsShow struct {
	Context
}

// Perform executes the action
func (a *ListingsShow) Perform() {
	id, err := strconv.Atoi(a.Param("id"))
	if err != nil {
		a.Diagnostics().AddError(err)
		a.RespondNotFound()
		return
	}

	publishedNull := map[string]interface{}{"published_at": nil}

	property := &flentities.Property{}

	err = a.Repository().Preload("Amenities").Preload("Images", flentities.Image{Uploaded: true}, flentities.ImagesDefaultOrder).
		Preload("Units", func(db *gorm.DB) *gorm.DB { return db.Not(publishedNull).Order("units.base_price_cents") }).
		Preload("Units.Images", flentities.Image{Uploaded: true}, flentities.ImagesDefaultOrder).Preload("Units.Amenities").
		Where(id).Find(property).Error

	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			a.RespondNotFound()
		} else {
			a.ServerError(err)
		}
		return
	}

	user := a.CurrentUser()
	if user == nil || user.AccountID == nil || *user.AccountID != property.AccountID {
		if property.PublishedAt == nil || len(property.Units) == 0 {
			a.RespondNotFound()
			return
		}
	}

	similarListings := []*flentities.Property{}
	err = a.Repository().Preload("Images", flentities.Image{Uploaded: true}, flentities.ImagesDefaultOrder,
		func(db *gorm.DB) *gorm.DB { return db.Limit(1) }).
		Preload("Units", func(db *gorm.DB) *gorm.DB {
			return db.Not(publishedNull).Order("units.base_price_cents").Limit(1)
		}).
		Where(map[string]interface{}{"city_id": property.CityID}).
		Where("id in (?)", a.Repository().Table("units").Select("property_id").Not(publishedNull).QueryExpr()).
		Not(publishedNull).Not(property.ID).Limit(3).Find(&similarListings).Error

	if err != nil {
		a.Diagnostics().Add("similar_listings_error", err.Error())
	}

	a.Respond(http.StatusOK, flviews.NewListing(property, similarListings))
}

// NewListingsShow returns a new ListingsShow action
func NewListingsShow(c Context) Action {
	return &ListingsShow{c}
}
