package fujilane

import (
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type amenityRow struct {
	flentities.Amenity
	Property string
	Unit     string
}

func requestAmenities(target string) error {
	return perform("GET", strings.Replace(flweb.AmenityTypesPath, ":target", target, 1), nil)
}

func tableRowToAmenity(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*amenityRow)
	return &row.Amenity, loadAssociationByName(&row.Amenity, "Property", row.Property, "Unit", row.Unit)
}

func amenityToTableRow(r *flentities.Repository, a interface{}) (interface{}, error) {
	amenity := a.(*flentities.Amenity)

	amenity.Property = &flentities.Property{}
	amenity.Unit = &flentities.Unit{}

	assocs := map[string]interface{}{
		"Property": amenity.Property,
		"Unit":     amenity.Unit,
	}

	if err := loadAssociations(r, a, assocs); err != nil {
		return nil, err
	}

	row := &amenityRow{Amenity: *amenity, Unit: amenity.Unit.Name}
	if amenity.Property.Name != nil {
		row.Property = *amenity.Property.Name
	}

	return row, nil
}

func AmenityContext(s *godog.Suite) {
	s.Step(`^the following amenities:$`, createFromTableStep(new(amenityRow), tableRowToAmenity))
	s.Step(`^I list amenity types for "([^"]*)"$`, requestAmenities)
	s.Step(`^the system should respond with "([^"]*)" and the following amenity types:$`,
		assertResponseStatusAndListStep(&[]*flentities.AmenityType{}))
	s.Step(`^I should have the following amenities:$`, assertDatabaseRecordsStep(&[]*flentities.Amenity{}, amenityToTableRow))
}
