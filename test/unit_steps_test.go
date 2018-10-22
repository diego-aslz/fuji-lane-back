package fujilane

import (
	"strconv"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type unitRow struct {
	flentities.Unit
	Property string
}

func requestUnitsCreate(table *gherkin.DataTable) error {
	unit, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	body := flactions.UnitsCreateBody{}
	body.Name = unit["Name"]

	body.Bedrooms, _ = strconv.Atoi(unit["Bedrooms"])
	body.SizeM2, _ = strconv.Atoi(unit["SizeM2"])
	body.MaxOccupancy, _ = strconv.Atoi(unit["MaxOccupancy"])
	body.Count, _ = strconv.Atoi(unit["Count"])

	propertyName := unit["Property"]
	if propertyName != "" {
		property := &flentities.Property{}
		if err := findByName(property, propertyName); err != nil {
			return err
		}
		body.PropertyID = property.ID
	}

	return performPOST(flweb.UnitsPath, body)
}

func unitToTableRow(r *flentities.Repository, a interface{}) (interface{}, error) {
	unit := a.(*flentities.Unit)

	unit.Property = &flentities.Property{}

	assocs := map[string]interface{}{
		"Property": unit.Property,
	}

	if err := loadAssociations(r, a, assocs); err != nil {
		return nil, err
	}

	row := &unitRow{Unit: *unit}
	if unit.Property.Name != nil {
		row.Property = *unit.Property.Name
	}

	return row, nil
}

func UnitContext(s *godog.Suite) {
	s.Step(`^I add the following unit:$`, requestUnitsCreate)
	s.Step(`^I should have the following units:$`, assertDatabaseRecordsStep(&[]*flentities.Unit{}, unitToTableRow))
	s.Step(`^I should have no units$`, assertNoDatabaseRecordsStep(&flentities.Unit{}))
}
