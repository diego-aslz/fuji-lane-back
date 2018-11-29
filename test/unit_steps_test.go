package fujilane

import (
	"fmt"
	"io"
	"strconv"
	"strings"

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

func requestUnitsShow(name string) error {
	unit := &flentities.Unit{}
	if err := findByName(unit, name); err != nil {
		return err
	}

	url := strings.Replace(flweb.UnitPath, ":id", fmt.Sprint(unit.ID), 1)

	return performGETStep(url)()
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

func requestUnitsUpdate(name string, table *gherkin.DataTable) error {
	u, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	unit := &flentities.Unit{}
	if err := findByName(unit, name); err != nil {
		return err
	}

	body := flactions.UnitsUpdateBody{}
	body.Name = u["Name"]

	body.Bedrooms, _ = strconv.Atoi(u["Bedrooms"])
	body.SizeM2, _ = strconv.Atoi(u["SizeM2"])
	body.MaxOccupancy, _ = strconv.Atoi(u["MaxOccupancy"])
	body.Count, _ = strconv.Atoi(u["Count"])
	body.BasePriceCents, _ = strconv.Atoi(u["BasePriceCents"])
	body.OneNightPriceCents, _ = strconv.Atoi(u["OneNightPriceCents"])
	body.OneWeekPriceCents, _ = strconv.Atoi(u["OneWeekPriceCents"])
	body.ThreeMonthsPriceCents, _ = strconv.Atoi(u["ThreeMonthsPriceCents"])
	body.SixMonthsPriceCents, _ = strconv.Atoi(u["SixMonthsPriceCents"])
	body.TwelveMonthsPriceCents, _ = strconv.Atoi(u["TwelveMonthsPriceCents"])

	FloorPlanImageID, _ := strconv.Atoi(u["FloorPlanImageID"])
	body.FloorPlanImageID = uint(FloorPlanImageID)

	var bodyIO io.Reader
	if bodyIO, err = bodyFromObject(body); err != nil {
		return err
	}

	return perform("PUT", strings.Replace(flweb.UnitPath, ":id", fmt.Sprint(unit.ID), 1), bodyIO)
}

func tableRowToUnit(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*unitRow)
	return &row.Unit, loadAssociationByName(&row.Unit, "Property", row.Property)
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

func updateUnit(name string, table *gherkin.DataTable) error {
	unit := &flentities.Unit{}
	if err := findByName(unit, name); err != nil {
		return err
	}

	tbl, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	updates := map[string]interface{}{}
	updates["FloorPlanImageID"], err = strconv.Atoi(tbl["FloorPlanImageID"])
	if err != nil {
		return err
	}

	return flentities.WithRepository(func(r *flentities.Repository) error {
		return r.Model(unit).Updates(updates).Error
	})
}

func requestUnitsUpdateWithAmenities(name string, table *gherkin.DataTable) error {
	b, err := assist.CreateSlice(new(flactions.AmenityBody), table)
	if err != nil {
		return err
	}

	unit := &flentities.Unit{}
	if err = findByName(unit, name); err != nil {
		return err
	}

	amenities := b.([]*flactions.AmenityBody)
	updateBody := &flactions.UnitsUpdateBody{}
	updateBody.Amenities = amenities
	body, err := bodyFromObject(updateBody)

	return perform("PUT", strings.Replace(flweb.UnitPath, ":id", fmt.Sprint(unit.ID), 1), body)
}

func requestUnitsPublish(id string) error {
	url := strings.Replace(flweb.UnitsPublishPath, ":id", id, 1)

	return perform("PUT", url, nil)
}

func UnitContext(s *godog.Suite) {
	s.Step(`^unit "([^"]*)" has:$`, updateUnit)
	s.Step(`^the following units:$`, createFromTableStep(new(unitRow), tableRowToUnit))
	s.Step(`^I add the following unit:$`, requestUnitsCreate)
	s.Step(`^I update unit "([^"]*)" with the following attributes:$`, requestUnitsUpdate)
	s.Step(`^I update unit "([^"]*)" with the following amenities:$`, requestUnitsUpdateWithAmenities)
	s.Step(`^I should have the following units:$`, assertDatabaseRecordsStep(&[]*flentities.Unit{}, unitToTableRow))
	s.Step(`^I should have no units$`, assertNoDatabaseRecordsStep(&flentities.Unit{}))
	s.Step(`^I get details for unit "([^"]*)"$`, requestUnitsShow)
	s.Step(`^I publish unit "([^"]*)"$`, requestUnitsPublish)
}
