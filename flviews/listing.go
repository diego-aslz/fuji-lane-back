package flviews

import (
	"github.com/nerde/fuji-lane-back/flentities"
)

// NewListing returns the expected structure for a Listing view
func NewListing(property *flentities.Property, similarListings []*flentities.Property) map[string]interface{} {
	return map[string]interface{}{
		"id":              property.ID,
		"name":            property.Name,
		"address1":        property.Address1,
		"address2":        property.Address2,
		"address3":        property.Address3,
		"cityID":          property.CityID,
		"postalCode":      property.PostalCode,
		"countryID":       property.CountryID,
		"latitude":        property.Latitude,
		"longitude":       property.Longitude,
		"overview":        property.Overview,
		"images":          listingImages(property.Images),
		"amenities":       listingAmenities(property.Amenities),
		"units":           listingUnits(property.Units),
		"similarListings": listingSimilarListings(similarListings),
	}
}

func listingUnits(units []*flentities.Unit) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, u := range units {
		result = append(result, map[string]interface{}{
			"id":             u.ID,
			"name":           u.Name,
			"bedrooms":       u.Bedrooms,
			"bathrooms":      u.Bathrooms,
			"sizeM2":         u.SizeM2,
			"maxOccupancy":   u.MaxOccupancy,
			"basePriceCents": u.BasePriceCents,
			"overview":       u.Overview,
			"amenities":      listingAmenities(u.Amenities),
			"images":         listingImages(u.Images),
		})
	}

	return result
}

func listingImages(images []*flentities.Image) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, img := range images {
		result = append(result, map[string]interface{}{
			"name": img.Name,
			"url":  img.URL,
		})
	}

	return result
}

func listingAmenities(amenities []*flentities.Amenity) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, am := range amenities {
		result = append(result, map[string]interface{}{
			"type": am.Type,
			"name": flentities.AmenityName(*am),
		})
	}

	return result
}

func listingSimilarListings(similarListings []*flentities.Property) []map[string]interface{} {
	result := []map[string]interface{}{}

	for _, l := range similarListings {
		if l.Units == nil || len(l.Units) == 0 {
			continue
		}

		result = append(result, map[string]interface{}{
			"name":           l.Name,
			"address1":       l.Address1,
			"address2":       l.Address2,
			"address3":       l.Address3,
			"overview":       l.Overview,
			"bedrooms":       l.Units[0].Bedrooms,
			"bathrooms":      l.Units[0].Bathrooms,
			"sizeM2":         l.Units[0].SizeM2,
			"basePriceCents": l.Units[0].BasePriceCents,
		})
	}

	return result
}
