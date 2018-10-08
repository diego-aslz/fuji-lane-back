package fujilane

import (
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

func requestPropertiesCreate() error {
	return performPOST(flweb.PropertiesPath, nil)
}

func requestPropertiesShow(name string) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		property := &flentities.Property{}
		if err := r.Find(property, map[string]interface{}{"name": name}).Error; err != nil {
			return err
		}

		url := strings.Replace(flweb.PropertiesShowPath, ":id", fmt.Sprint(property.ID), 1)

		return performGET(url)
	})
}

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

func (row *propertyRow) save(r *flentities.Repository) error {
	if row.Account != "" {
		row.Property.Account = &flentities.Account{}
		err := r.Find(row.Property.Account, flentities.Account{Name: row.Account}).Error
		if err != nil {
			return err
		}
	}

	if row.City != "" {
		row.Property.City = &flentities.City{}
		err := r.Find(row.Property.City, flentities.City{Name: row.City}).Error
		if err != nil {
			return err
		}
	}

	if row.Country != "" {
		row.Property.Country = &flentities.Country{}
		err := r.Find(row.Property.Country, flentities.Country{Name: row.Country}).Error
		if err != nil {
			return err
		}
	}

	if row.Name != "" {
		row.Property.Name = &row.Name
	}

	if row.Address1 != "" {
		row.Property.Address1 = &row.Address1
	}

	if row.Address2 != "" {
		row.Property.Address2 = &row.Address2
	}

	if row.Address3 != "" {
		row.Property.Address3 = &row.Address3
	}

	if row.PostalCode != "" {
		row.Property.PostalCode = &row.PostalCode
	}

	switch row.State {
	default:
		row.Property.StateID = flentities.PropertyStateDraft
	}

	return r.Create(&row.Property).Error
}

func theFollowingProperties(table *gherkin.DataTable) error {
	return createRowEntitiesFromTable(new(propertyRow), table)
}

func assertProperties(table *gherkin.DataTable) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		properties := []*flentities.Property{}
		err := r.Preload("Account").Find(&properties).Error
		if err != nil {
			return err
		}

		rowsCount := len(table.Rows) - 1
		if len(properties) != rowsCount {
			return fmt.Errorf("Expected to have %d properties in the DB, got %d", rowsCount, len(properties))
		}

		rows := []*propertyRow{}
		for _, p := range properties {
			row := &propertyRow{Property: *p, State: p.State()}
			if p.Account != nil {
				row.Account = p.Account.Name
			}

			rows = append(rows, row)
		}

		return assist.CompareToSlice(rows, table)
	})
}

func assertNoProperties() error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		count := 0
		err := r.Model(&flentities.Property{}).Count(&count).Error
		if err != nil {
			return err
		}

		if count != 0 {
			return fmt.Errorf("Expected to have %d properties in the DB, got %d", 0, count)
		}

		return nil
	})
}

func PropertyContext(s *godog.Suite) {
	s.Step(`^I add a new property$`, requestPropertiesCreate)
	s.Step(`^the following properties:$`, theFollowingProperties)
	s.Step(`^we should have the following properties:$`, assertProperties)
	s.Step(`^we should have no properties$`, assertNoProperties)
	s.Step(`^I get details for property "([^"]*)"$`, requestPropertiesShow)
}
