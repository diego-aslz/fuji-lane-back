package fujilane

import (
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
)

func theFollowingCountries(table *gherkin.DataTable) error {
	return createEntitiesFromTable(new(flentities.Country), table)
}

func assertCountries(table *gherkin.DataTable) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		count := 0
		err := r.Model(&flentities.Country{}).Count(&count).Error
		if err != nil {
			return err
		}

		rowsCount := len(table.Rows) - 1
		if count != rowsCount {
			return fmt.Errorf("Expected to have %d countries in the DB, got %d", rowsCount, count)
		}

		countries := []*flentities.Country{}
		err = r.Find(&countries).Error
		if err != nil {
			return err
		}

		return assist.CompareToSlice(countries, table)
	})
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^the following countries:$`, theFollowingCountries)
	s.Step(`^we should have the following countries:$`, assertCountries)
}
