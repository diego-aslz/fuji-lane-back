package flactions

import "github.com/nerde/fuji-lane-back/flentities"

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

type bodyWithAmenities struct {
	Amenities []*AmenityBody `json:"amenities"`
}

func (a bodyWithAmenities) amenitiesDiff(previousAmenities []*flentities.Amenity) (
	amenitiesToDelete []*flentities.Amenity, amenitiesToCreate []*flentities.Amenity) {

	if a.Amenities == nil {
		return
	}

	// Checking which amenities were removed by the user so we can delete them from the database
	for _, am := range previousAmenities {
		removedByUser := true
		for _, ab := range a.Amenities {
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

	for _, ab := range a.Amenities {
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

		am := &flentities.Amenity{Type: ab.Type, Name: &ab.Name}
		if ab.Type != "custom" {
			am.Name = nil
		}

		amenitiesToCreate = append(amenitiesToCreate, am)
	}

	return
}
