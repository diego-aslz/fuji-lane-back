package flentities

import (
	"sort"
	"strconv"
	"sync"

	"github.com/jinzhu/gorm"
)

// ListingsSearchFilters are the filters to apply to a search
type ListingsSearchFilters struct {
	CityID            uint
	MinBedrooms       int
	MinBathrooms      int
	Page              int
	PerPage           int
	CheckIn           *Date
	CheckOut          *Date
	MinPriceCents     int
	MaxPriceCents     int
	UnitAmenities     []string
	PropertyAmenities []string
	nights            int
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

func (f *ListingsSearchFilters) addUnitAmenityConditions(db *gorm.DB) *gorm.DB {
	return addAmenityConditions(db, f.UnitAmenities, UnitAmenityTypes, "units", "unit_id")
}

func (f *ListingsSearchFilters) addPropertyAmenityConditions(db *gorm.DB) *gorm.DB {
	return addAmenityConditions(db, f.PropertyAmenities, PropertyAmenityTypes, "properties", "property_id")
}

func addAmenityConditions(db *gorm.DB, amenities []string, types []*AmenityType, joinTable, fk string) *gorm.DB {
	for i, aType := range amenities {
		found := false
		for _, a := range types {
			if a.Code == aType {
				found = true
				break
			}
		}

		if !found {
			continue
		}

		alias := joinTable + "_amenities_" + strconv.Itoa(i)
		join := "INNER JOIN amenities " + alias + " ON " + joinTable + ".id = " + alias + "." + fk
		join += " AND " + alias + ".type = ?"
		db = db.Joins(join, aType)
	}

	return db
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

	unitConditions = s.addUnitAmenityConditions(unitConditions)

	unitConditions = unitConditions.Joins("INNER JOIN properties ON units.property_id = properties.id")
	unitConditions = s.addPropertyAmenityConditions(unitConditions)

	joinedUnits := unitConditions.
		Where("properties.city_id = ?", s.CityID).
		Where("properties.published_at IS NOT NULL")

	if s.hasDates() {
		joinedUnits = joinedUnits.Where("properties.minimum_stay <= ?", s.Nights())
	}

	result := &ListingsSearchResult{Properties: []*Property{}}

	var routineErr error
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func(joinedUnits *gorm.DB) {
		defer wg.Done()
		routineErr = s.addMetadata(joinedUnits, result)
	}(joinedUnits)

	join, args := s.pricesJoin()
	unitConditions = unitConditions.Joins(join, args...)

	if s.MinPriceCents > 0 || s.MaxPriceCents > 0 {
		joinedUnits = joinedUnits.Joins(join, args...)
	}

	propertyConditions := s.
		Where("id IN (?)", joinedUnits.Model(&Unit{}).Select("units.property_id").QueryExpr()).
		Preload("Images", Image{Uploaded: true}, ImagesDefaultOrder).
		Preload("Amenities")

	propertyConditions = Repository{propertyConditions}.Paginate(s.Page, s.PerPage)

	var err error
	if err = propertyConditions.Find(&result.Properties).Error; err != nil {
		wg.Wait()
		return result, err
	}

	units := []*Unit{}
	propertyIDs := []uint{}
	for _, p := range result.Properties {
		propertyIDs = append(propertyIDs, p.ID)
	}

	err = unitConditions.
		Select("DISTINCT units.*").
		Where("units.property_id IN (?)", propertyIDs).
		Preload("Images", Image{Uploaded: true}, ImagesDefaultOrder).
		Preload("Amenities").
		Preload("Prices", func(db *gorm.DB) *gorm.DB {
			if s.hasDates() {
				db = db.Where("min_nights <= ?", s.Nights())
			}

			return db.Order("min_nights")
		}).
		Find(&units).
		Error

	if err != nil {
		wg.Wait()
		return result, err
	}

	for _, p := range result.Properties {
		p.Units = []*Unit{}
		for _, u := range units {
			if u.PropertyID == p.ID {
				p.Units = append(p.Units, u)
			}
		}

		sort.Slice(p.Units, func(i, j int) bool {
			li := len(p.Units[i].Prices)
			lj := len(p.Units[j].Prices)

			return p.Units[i].Prices[li-1].PerNightCents() < p.Units[j].Prices[lj-1].PerNightCents()
		})
	}

	wg.Wait()

	if err == nil {
		err = routineErr
	}

	return result, err
}

func (s ListingsSearch) addMetadata(unitConditions *gorm.DB, result *ListingsSearchResult) error {
	bestPriceJoin := "INNER JOIN (" +
		"SELECT unit_id, MAX(min_nights) min_nights FROM prices WHERE min_nights <= ? GROUP BY unit_id" +
		") best_price ON best_price.unit_id = units.id"

	selects := "COUNT(distinct units.property_id)"
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
