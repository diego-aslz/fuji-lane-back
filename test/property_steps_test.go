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
	return flentities.WithRepository(func(r *flentities.Repository) error {
		property := &flentities.Property{}
		if err := r.Find(property, map[string]interface{}{"name": name}).Error; err != nil {
			return err
		}

		url := strings.Replace(flweb.PropertyPath, ":id", fmt.Sprint(property.ID), 1)

		return performGETStep(url)()
	})
}

func tableRowToProperty(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*propertyRow)

	if row.Account != "" {
		row.Property.Account = &flentities.Account{}
		err := r.Find(row.Property.Account, flentities.Account{Name: row.Account}).Error
		if err != nil {
			return nil, err
		}
	}

	if row.City != "" {
		row.Property.City = &flentities.City{}
		err := r.Find(row.Property.City, flentities.City{Name: row.City}).Error
		if err != nil {
			return nil, err
		}
	}

	if row.Country != "" {
		row.Property.Country = &flentities.Country{}
		err := r.Find(row.Property.Country, flentities.Country{Name: row.Country}).Error
		if err != nil {
			return nil, err
		}
	}

	row.Property.Name = refStr(row.Name)
	row.Property.Address1 = refStr(row.Address1)
	row.Property.Address2 = refStr(row.Address2)
	row.Property.Address3 = refStr(row.Address3)
	row.Property.PostalCode = refStr(row.PostalCode)

	switch row.State {
	default:
		row.Property.StateID = flentities.PropertyStateDraft
	}

	return &row.Property, nil
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
