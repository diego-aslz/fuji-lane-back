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

func tableRowToCity(r *flentities.Repository, a interface{}) (interface{}, error) {
	row := a.(*cityRow)

	if row.Country != "" {
		err := r.Find(&row.City.Country, flentities.Country{Name: row.Country}).Error
		if err != nil {
			return nil, err
		}
	}

	return &row.City, nil
}

func requestCities() error {
	return performGET(flweb.CitiesPath)
}

func assertCities(table *gherkin.DataTable) error {
	return assertDatabaseRecords(&[]*flentities.City{}, table)
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
	s.Step(`^the following cities:$`, createFromTableStep(new(cityRow), tableRowToCity))
	s.Step(`^I list cities$`, requestCities)
	s.Step(`^we should have the following cities:$`, assertCities)
	s.Step(`^the system should respond with "([^"]*)" and the following cities:$`, assertCitiesResponse)
}
