package fujilane

import (
	"encoding/json"
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

type cityRow struct {
	flentities.City
	Country string
}

func (row *cityRow) save(r *flentities.Repository) error {
	if row.Country != "" {
		err := r.Find(&row.City.Country, flentities.Country{Name: row.Country}).Error
		if err != nil {
			return err
		}
	}
	return r.Create(&row.City).Error
}

func theFollowingCities(table *gherkin.DataTable) error {
	return createRowEntitiesFromTable(new(cityRow), table)
}

func requestCities() error {
	return performGET(flweb.CitiesPath)
}

func assertCities(table *gherkin.DataTable) error {
	return flentities.WithRepository(func(r *flentities.Repository) error {
		count := 0
		err := r.Model(&flentities.City{}).Count(&count).Error
		if err != nil {
			return err
		}

		rowsCount := len(table.Rows) - 1
		if count != rowsCount {
			return fmt.Errorf("Expected to have %d cities in the DB, got %d", rowsCount, count)
		}

		cities := []*flentities.City{}
		err = r.Find(&cities).Error
		if err != nil {
			return err
		}

		return assist.CompareToSlice(cities, table)
	})
}

func assertCitiesResponse(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	cities := []*flentities.City{}
	if err := json.Unmarshal([]byte(response.Body.String()), &cities); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	return assist.CompareToSlice(cities, table)
}

func CityContext(s *godog.Suite) {
	s.Step(`^the following cities:$`, theFollowingCities)
	s.Step(`^I list cities$`, requestCities)
	s.Step(`^we should have the following cities:$`, assertCities)
	s.Step(`^the system should respond with "([^"]*)" and the following cities:$`, assertCitiesResponse)
}
