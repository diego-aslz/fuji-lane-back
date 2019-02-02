package flentities

import (
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
		join += " AND " + PerNightPriceSQL + " >= ?"
		args = append(args, f.MinPriceCents)
	}

	if f.MaxPriceCents > 0 {
		join += " AND " + PerNightPriceSQL + " <= ?"
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
func (s ListingsSearch) Search() (*ListingsSearchResult, error) {
	unitConditions := s.Where("units.published_at IS NOT NULL")

	if s.MinBedrooms > 0 {
		unitConditions = unitConditions.Where("bedrooms >= ?", s.MinBedrooms)
	}

	if s.MinBathrooms > 0 {
		unitConditions = unitConditions.Where("bathrooms >= ?", s.MinBathrooms)
	}

	joinedUnits := unitConditions.
		Joins("INNER JOIN properties ON units.property_id = properties.id").
		Where("properties.city_id = ?", s.CityID).
		Where("properties.published_at IS NOT NULL")

	if s.hasDates() {
		joinedUnits = joinedUnits.Where("properties.minimum_stay <= ?", s.Nights())
	}

	result := &ListingsSearchResult{Properties: []*Property{}}

	if err := s.addMetadata(joinedUnits, result); err != nil {
		return result, err
	}

	join, args := s.pricesJoin()
	unitConditions = unitConditions.Joins(join, args...)

	if s.MinPriceCents > 0 || s.MaxPriceCents > 0 {
		joinedUnits = joinedUnits.Joins(join, args...)
	}

	propertyConditions := s.
		Where("id IN (?)", joinedUnits.Model(&Unit{}).Select("property_id").QueryExpr()).
		Preload("Images", Image{Uploaded: true}, ImagesDefaultOrder).
		Preload("Amenities").
		Preload("Units", func(_ *gorm.DB) *gorm.DB { return unitConditions.Order("prices.cents / prices.min_nights") }).
		Preload("Units.Images", Image{Uploaded: true}, ImagesDefaultOrder).
		Preload("Units.Amenities").
		Preload("Units.Prices")

	propertyConditions = Repository{propertyConditions}.Paginate(s.Page, s.PerPage)

	return result, propertyConditions.Find(&result.Properties).Error
}

func (s ListingsSearch) addMetadata(unitConditions *gorm.DB, result *ListingsSearchResult) error {
	bestPriceJoin := "INNER JOIN (" +
		"SELECT unit_id, MAX(min_nights) min_nights FROM prices WHERE min_nights <= ? GROUP BY unit_id" +
		") best_price ON best_price.unit_id = units.id"

	selects := "COUNT(distinct property_id)"
	selects += ", COALESCE(ROUND(MIN(" + PerNightPriceSQL + ")), 0)"
	selects += ", COALESCE(ROUND(MAX(" + PerNightPriceSQL + ")), 0)"
	selects += ", COALESCE(ROUND(AVG(" + PerNightPriceSQL + ")), 0)"

	rows, err := unitConditions.
		Model(&Unit{}).
		Joins(bestPriceJoin, s.Nights()).
		Joins("INNER JOIN prices ON prices.unit_id = units.id AND prices.min_nights = best_price.min_nights").
		Select(selects).
		Rows()

	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&result.TotalPropertiesCount, &result.MinPerNightCents, &result.MaxPerNightCents,
			&result.AvgPerNightCents)
	}

	return err
}

// ListingsSearchResult represents the result of a search performed
type ListingsSearchResult struct {
	Properties           []*Property
	TotalPropertiesCount int
	MinPerNightCents     int
	MaxPerNightCents     int
	AvgPerNightCents     int
}
