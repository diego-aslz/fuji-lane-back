package fujilane

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flactions"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type propertyRow struct {
	flentities.Property
	Account string
	City    string
	Country string
}

func requestPropertiesCreate() error {
	return performPOST(flweb.PropertiesPath, nil)
}

func requestPropertiesUpdate(id string, table *gherkin.DataTable) error {
	b, err := assist.CreateInstance(new(flactions.PropertiesUpdateBody), table)
	if err != nil {
		return err
	}

	body, err := bodyFromObject(b)

	return perform("PUT", strings.Replace(flweb.PropertyPath, ":id", id, 1), body)
}

func requestPropertiesUpdateWithAmenities(id string, table *gherkin.DataTable) error {
	b, err := assist.CreateSlice(new(flactions.AmenityBody), table)
	if err != nil {
		return err
	}

	amenities := b.([]*flactions.AmenityBody)
	updateBody := &flactions.PropertiesUpdateBody{}
	updateBody.Amenities = amenities
	body, err := bodyFromObject(updateBody)

	return perform("PUT", strings.Replace(flweb.PropertyPath, ":id", id, 1), body)
}

func requestPropertiesShow(name string) error {
	property := &flentities.Property{}
	if err := findByName(property, name); err != nil {
		return err
	}

	url := strings.Replace(flweb.PropertyPath, ":id", fmt.Sprint(property.ID), 1)

	return performGETStep(url)()
}

func tableRowToProperty(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*propertyRow)

	return &row.Property, loadAssociationByName(&row.Property,
		"Account", row.Account,
		"City", row.City,
		"Country", row.Country)
}

func propertyToTableRow(r *flentities.Repository, p interface{}) (interface{}, error) {
	property := p.(*flentities.Property)

	property.Account = &flentities.Account{}
	property.City = &flentities.City{}
	property.Country = &flentities.Country{}

	assocs := map[string]interface{}{
		"Account": property.Account,
		"City":    property.City,
		"Country": property.Country,
	}

	if err := loadAssociations(r, p, assocs); err != nil {
		return nil, err
	}

	row := &propertyRow{
		Property: *property,
		Account:  property.Account.Name,
		City:     property.City.Name,
		Country:  property.Country.Name,
	}

	return row, nil
}

func PropertyContext(s *godog.Suite) {
	s.Step(`^I add a new property$`, requestPropertiesCreate)
	s.Step(`^I update the property "([^"]*)" with the following details:$`, requestPropertiesUpdate)
	s.Step(`^I update the property "([^"]*)" with the following amenities:$`, requestPropertiesUpdateWithAmenities)
	s.Step(`^the following properties:$`, createFromTableStep(new(propertyRow), tableRowToProperty))
	s.Step(`^I should have the following properties:$`, assertDatabaseRecordsStep(&[]*flentities.Property{}, propertyToTableRow))
	s.Step(`^I should have no properties$`, assertNoDatabaseRecordsStep(&flentities.Property{}))
	s.Step(`^I get details for property "([^"]*)"$`, requestPropertiesShow)
}
