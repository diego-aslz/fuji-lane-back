package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/jinzhu/gorm"
)

func simulateAddProperty() error {
	return makePOSTRequest(propertiesPath, &gherkin.DataTable{})
}

type propertyRow struct {
	*Property
	User  string
	State string
}

func assertProperties(table *gherkin.DataTable) error {
	return withDatabase(func(db *gorm.DB) error {
		properties := []*Property{}
		err := db.Preload("User").Find(&properties).Error
		if err != nil {
			return err
		}

		rowsCount := len(table.Rows) - 1
		if len(properties) != rowsCount {
			return fmt.Errorf("Expected to have %d properties in the DB, got %d", rowsCount, len(properties))
		}

		rows := []*propertyRow{}
		for _, p := range properties {
			userName := ""
			if p.User != nil {
				userName = p.User.Name
			}

			rows = append(rows, &propertyRow{p, userName, p.State()})
		}

		return assist.CompareToSlice(rows, table)
	})
}

func assertNoProperties() error {
	return withDatabase(func(db *gorm.DB) error {
		count := 0
		err := db.Model(&Property{}).Count(&count).Error
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
	s.Step(`^I should have the following properties:$`, assertProperties)
	s.Step(`^I should have no properties$`, assertNoProperties)
}
