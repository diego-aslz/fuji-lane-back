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
	return &row.City, loadAssociationByName(&row.City, "Country", row.Country)
}

func CityContext(s *godog.Suite) {
	s.Step(`^the following cities:$`, createFromTableStep(new(cityRow), tableRowToCity))
	s.Step(`^I list cities$`, performGETStep(flweb.CitiesPath))
	s.Step(`^I should have the following cities:$`, assertDatabaseRecordsStep(&[]*flentities.City{}))
	s.Step(`^I should receive an "([^"]*)" response with the following cities:$`,
		assertResponseStatusAndListStep(&[]*flentities.City{}))
}
