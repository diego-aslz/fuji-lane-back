package flentities

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

// ListingsSearchFilters are the filters to apply to a search
type ListingsSearchFilters struct {
	CityID       uint
	MinBedrooms  int
	MinBathrooms int
	Page         int
	PerPage      int
}

// ListingsSearch for searching for listings
type ListingsSearch struct {
	*Repository
	ListingsSearchFilters
}

// Search searches for properties matching the filters
func (ps ListingsSearch) Search() ([]*Property, error) {
	publishedNull := map[string]interface{}{"published_at": nil}
	unitConditions := ps.Not(publishedNull).Where(map[string]interface{}{"deleted_at": nil})
	unitRawConditions := []string{"units.published_at IS NOT NULL", "units.deleted_at IS NULL"}
	unitJoinArgs := []interface{}{}

	if ps.MinBedrooms > 0 {
		condition := "bedrooms >= ?"
		unitConditions = unitConditions.Where(condition, ps.MinBedrooms)
		unitRawConditions = append(unitRawConditions, condition)
		unitJoinArgs = append(unitJoinArgs, ps.MinBedrooms)
	}

	if ps.MinBathrooms > 0 {
		condition := "bathrooms >= ?"
		unitConditions = unitConditions.Where(condition, ps.MinBathrooms)
		unitRawConditions = append(unitRawConditions, condition)
		unitJoinArgs = append(unitJoinArgs, ps.MinBathrooms)
	}

	builder := ps.
		Preload("Images", Image{Uploaded: true}, ImagesDefaultOrder).
		Preload("Amenities").
		Preload("Units", func(_ *gorm.DB) *gorm.DB { return unitConditions.Order("base_price_cents") }).
		Preload("Units.Images", Image{Uploaded: true}, ImagesDefaultOrder).
		Preload("Units.Amenities").
		Where("city_id = ?", ps.CityID).
		Joins(fmt.Sprintf("INNER JOIN units ON properties.id = units.property_id AND %s",
			strings.Join(unitRawConditions, " AND ")), unitJoinArgs...).
		Select("DISTINCT(properties.*)")

	properties := []*Property{}

	return properties, Repository{builder}.Paginate(ps.Page, ps.PerPage).Find(&properties).Error
}
