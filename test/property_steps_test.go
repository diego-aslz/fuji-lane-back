package fujilane

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gin-gonic/gin"
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

func requestPropertiesCreate(table *gherkin.DataTable) error {
	b, err := assist.CreateInstance(new(flactions.PropertiesCreateBody), table)

	if err != nil {
		return err
	}

	return performPOST(flweb.PropertiesPath, b)
}

func requestPropertiesUpdate(id string, table *gherkin.DataTable) error {
	b, err := assist.CreateInstance(new(flactions.PropertiesUpdateBody), table)
	if err != nil {
		return err
	}

	body, err := bodyFromObject(slicePayload(b, tableColumn(table, 0)))

	return perform("PUT", strings.Replace(flweb.PropertyPath, ":id", id, 1), body)
}

func requestPropertiesUpdateWithAmenities(id string, table *gherkin.DataTable) error {
	b, err := assist.CreateSlice(new(flactions.AmenityBody), table)
	if err != nil {
		return err
	}

	amenities := b.([]*flactions.AmenityBody)
	updateBody := gin.H{"amenities": amenities}
	body, err := bodyFromObject(updateBody)

	return perform("PUT", strings.Replace(flweb.PropertyPath, ":id", id, 1), body)
}

func performGETPropertyStep(path string) func(string) error {
	return func(propertyName string) error {
		property := &flentities.Property{}
		if err := findByName(property, propertyName); err != nil {
			return err
		}

		url := strings.Replace(path, ":id", fmt.Sprint(property.ID), 1)

		return performGETStep(url)()
	}
}

func requestPropertiesPublish(id string) error {
	url := strings.Replace(flweb.PropertiesPublishPath, ":id", id, 1)

	return perform("PUT", url, nil)
}

func requestPropertiesUnpublish(id string) error {
	url := strings.Replace(flweb.PropertiesUnpublishPath, ":id", id, 1)

	return perform("PUT", url, nil)
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
	s.Step(`^I create the following property:$`, requestPropertiesCreate)
	s.Step(`^I update the property "([^"]*)" with the following details:$`, requestPropertiesUpdate)
	s.Step(`^I update the property "([^"]*)" with the following amenities:$`, requestPropertiesUpdateWithAmenities)
	s.Step(`^I publish property "([^"]*)"$`, requestPropertiesPublish)
	s.Step(`^I unpublish property "([^"]*)"$`, requestPropertiesUnpublish)
	s.Step(`^the following properties:$`, createFromTableStep(new(propertyRow), tableRowToProperty))
	s.Step(`^I should have the following properties:$`, assertDatabaseRecordsStep(&[]*flentities.Property{}, propertyToTableRow))
	s.Step(`^I should have no properties$`, assertNoDatabaseRecordsStep(&flentities.Property{}))
	s.Step(`^I get details for property "([^"]*)"$`, performGETPropertyStep(flweb.PropertyPath))
	s.Step(`^I get listing details for "([^"]*)"$`, performGETPropertyStep(flweb.ListingPath))
	s.Step(`^I list my properties$`, performGETStep(flweb.PropertiesPath))
	s.Step(`^I get properties sitemap$`, performGETStep(flweb.PropertiesSitemapPath))
}
