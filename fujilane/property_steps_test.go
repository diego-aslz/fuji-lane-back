package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
)

func simulateAddProperty() error {
	return performPOSTWithTable(propertiesPath, &gherkin.DataTable{})
}

type propertyRow struct {
	*flentities.Property
	Account string
	State   string
}

func assertProperties(table *gherkin.DataTable) error {
	return withRepository(func(r *flentities.Repository) error {
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
			row := &propertyRow{Property: p, State: p.State()}
			if p.Account != nil {
				row.Account = p.Account.Name
			}

			rows = append(rows, row)
		}

		return assist.CompareToSlice(rows, table)
	})
}

func assertNoProperties() error {
	return withRepository(func(r *flentities.Repository) error {
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
	s.Step(`^I add a new property$`, simulateAddProperty)
	s.Step(`^we should have the following properties:$`, assertProperties)
	s.Step(`^we should have no properties$`, assertNoProperties)
}
