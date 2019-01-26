package flentities

import (
	"strings"

	"github.com/jinzhu/gorm"
)

// ListingsSearchFilters are the filters to apply to a search
type ListingsSearchFilters struct {
	CityID        uint
	MinBedrooms   int
	MinBathrooms  int
	Page          int
	PerPage       int
	CheckIn       *Date
	CheckOut      *Date
	MinPriceCents int
	MaxPriceCents int
	nights        int
}

// Nights returns how many nights is this query for
func (f *ListingsSearchFilters) Nights() int {
	if f.nights == 0 {
		if f.hasDates() {
			f.nights = f.CheckOut.NightsFrom(*f.CheckIn)
		} else {
			f.nights = 1
		}
	}

	return f.nights
}

func (f ListingsSearchFilters) hasDates() bool {
	return f.CheckIn != nil && f.CheckOut != nil
}

func (f *ListingsSearchFilters) pricesJoin() (string, []interface{}) {
	join := "INNER JOIN prices ON units.id = prices.unit_id"
	args := []interface{}{}

	// Apply nights even when there are no check in and check out dates available because we'll show the user the
	// price for one single night in that scenario. This avoids confusion if the user applies a price range filter and
	// sees listings which for long stays would match the max price limit, but we'd still show the price for
	// a single night to them, which would be above the max price limit.
	join += " AND prices.min_nights <= ?"
	args = append(args, f.Nights())

	if f.MinPriceCents > 0 {
		join += " AND prices.cents / prices.min_nights >= ?"
		args = append(args, f.MinPriceCents)
	}

	if f.MaxPriceCents > 0 {
		join += " AND prices.cents / prices.min_nights <= ?"
		args = append(args, f.MaxPriceCents)
	}

	return join, args
}

// ListingsSearch for searching for listings
type ListingsSearch struct {
	*Repository
	*ListingsSearchFilters
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
		Preload("Units", func(_ *gorm.DB) *gorm.DB {
			join, args := ps.pricesJoin()
			unitConditions = unitConditions.Joins(join, args...)

			return unitConditions.Order("prices.cents / prices.min_nights")
		}).
		Preload("Units.Images", Image{Uploaded: true}, ImagesDefaultOrder).
		Preload("Units.Amenities").
		Preload("Units.Prices").
		Where("city_id = ?", ps.CityID).
		Joins("INNER JOIN units ON properties.id = units.property_id AND "+strings.Join(unitRawConditions, " AND "),
			unitJoinArgs...).
		Select("DISTINCT(properties.*)")

	if ps.hasDates() {
		builder = builder.Where("minimum_stay <= ?", ps.Nights())
	}

	if ps.MinPriceCents > 0 || ps.MaxPriceCents > 0 {
		join, args := ps.pricesJoin()
		builder = builder.Joins(join, args...)
	}

	properties := []*Property{}

	return properties, Repository{builder}.Paginate(ps.Page, ps.PerPage).Find(&properties).Error
}
