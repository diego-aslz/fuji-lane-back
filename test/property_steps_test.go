package fujilane

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/jinzhu/gorm"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type propertyRow struct {
	flentities.Property
	Account    string
	State      string
	Name       string
	Address1   string
	Address2   string
	Address3   string
	PostalCode string
	City       string
	Country    string
}

func requestPropertiesCreate() error {
	return performPOST(flweb.PropertiesPath, nil)
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
	row.Property.Name = refStr(row.Name)
	row.Property.Address1 = refStr(row.Address1)
	row.Property.Address2 = refStr(row.Address2)
	row.Property.Address3 = refStr(row.Address3)
	row.Property.PostalCode = refStr(row.PostalCode)

	switch row.State {
	default:
		row.Property.StateID = flentities.PropertyStateDraft
	}

	return &row.Property, loadAssociationByName(&row.Property,
		"Account", row.Account,
		"City", row.City,
		"Country", row.Country)
}

func propertyToTableRow(r *flentities.Repository, p interface{}) (interface{}, error) {
	property := p.(*flentities.Property)

	property.Account = &flentities.Account{}
	err := r.Model(p).Association("Account").Find(property.Account).Error

	if err != nil && !gorm.IsRecordNotFoundError(err) {
		err = nil
	}

	return &propertyRow{Property: *property, State: property.State(), Account: property.Account.Name}, err
}

func PropertyContext(s *godog.Suite) {
	s.Step(`^I add a new property$`, requestPropertiesCreate)
	s.Step(`^the following properties:$`, createFromTableStep(new(propertyRow), tableRowToProperty))
	s.Step(`^we should have the following properties:$`, assertDatabaseRecordsStep(&[]*flentities.Property{}, propertyToTableRow))
	s.Step(`^we should have no properties$`, assertNoDatabaseRecordsStep(&flentities.Property{}))
	s.Step(`^I get details for property "([^"]*)"$`, requestPropertiesShow)
}
