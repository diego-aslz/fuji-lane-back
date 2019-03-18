package fujilane

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
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

	row := &amenityRow{Amenity: *amenity, Unit: amenity.Unit.Name, Property: amenity.Property.Name}

	return row, nil
}

func assertAmenityTypes(status, count string, table *gherkin.DataTable) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	types := []*flentities.AmenityType{}

	if err := json.Unmarshal([]byte(response.Body.String()), &types); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	c, _ := strconv.Atoi(count)

	if c != len(types) {
		return fmt.Errorf("Expected %d amenity types, got %d", c, len(types))
	}

	sample, err := assist.ParseSlice(table)
	if err != nil {
		return err
	}

	for idx := len(types) - 1; idx >= 0; idx-- {
		typ := types[idx]

		found := false
		for _, t := range sample {
			if t["Code"] == typ.Code && t["Name"] == typ.Name {
				t["Found"] = "t"
				found = true
				break
			}
		}

		if !found {
			types = append(types[:idx], types[idx+1:]...)
		}
	}

	missing := []string{}
	for _, t := range sample {
		if _, ok := t["Found"]; !ok {
			missing = append(missing, fmt.Sprint(t))
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("Missing amenity types:\n%s", strings.Join(missing, "\n"))
	}

	return nil
}

func AmenityContext(s *godog.Suite) {
	s.Step(`^the following amenities:$`, createFromTableStep(new(amenityRow), tableRowToAmenity))
	s.Step(`^I list amenity types for "([^"]*)"$`, requestAmenities)
	s.Step(`^I should receive an "([^"]*)" response with the following amenity types:$`,
		assertResponseStatusAndListStep(&[]*flentities.AmenityType{}))
	s.Step(`^I should receive an "([^"]*)" response with "([^"]*)" amenity types like:$`, assertAmenityTypes)
	s.Step(`^I should have the following amenities:$`, assertDatabaseRecordsStep(&[]*flentities.Amenity{}, amenityToTableRow))
}
