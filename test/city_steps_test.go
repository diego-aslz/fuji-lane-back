package fujilane

import (
	"github.com/DATA-DOG/godog"
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

func CityContext(s *godog.Suite) {
	s.Step(`^the following cities:$`, createFromTableStep(new(cityRow), tableRowToCity))
	s.Step(`^I list cities$`, performGETStep(flweb.CitiesPath))
	s.Step(`^we should have the following cities:$`, assertDatabaseRecordsStep(&[]*flentities.City{}))
	s.Step(`^the system should respond with "([^"]*)" and the following cities:$`,
		assertResponseStatusAndListStep(&[]*flentities.City{}))
}
