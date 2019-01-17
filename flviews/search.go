package flviews

import (
	"github.com/nerde/fuji-lane-back/flentities"
)

// NewSearch returns an array of maps that expose details about properties for Search endpoint
func NewSearch(properties []*flentities.Property) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, p := range properties {
		result = append(result, searchProperty(p))
	}

	return result
}

func searchProperty(property *flentities.Property) map[string]interface{} {
	return map[string]interface{}{
		"id":         property.ID,
		"name":       property.Name,
		"address1":   property.Address1,
		"address2":   property.Address2,
		"address3":   property.Address3,
		"cityID":     property.CityID,
		"postalCode": property.PostalCode,
		"countryID":  property.CountryID,
		"latitude":   property.Latitude,
		"longitude":  property.Longitude,
		"overview":   property.Overview,
		"images":     listingImages(property.Images),
		"amenities":  listingAmenities(property.Amenities),
		"units":      searchUnits(property.Units),
	}
}

func searchUnits(units []*flentities.Unit) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, u := range units {
		result = append(result, searchUnit(u))
	}

	return result
}

func searchUnit(u *flentities.Unit) map[string]interface{} {
	return map[string]interface{}{
		"id":                 u.ID,
		"name":               u.Name,
		"bedrooms":           u.Bedrooms,
		"bathrooms":          u.Bathrooms,
		"sizeM2":             u.SizeM2,
		"maxOccupancy":       u.MaxOccupancy,
		"perNightPriceCents": u.BasePriceCents,
		"amenities":          listingAmenities(u.Amenities),
		"images":             listingImages(u.Images),
	}
}