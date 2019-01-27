package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

func CountryContext(s *godog.Suite) {
	s.Step(`^the following countries:$`, createFromTableStep(new(flentities.Country)))
	s.Step(`^I should have the following countries:$`, assertDatabaseRecordsStep(&[]*flentities.Country{}))
	s.Step(`^I list countries$`, performGETStep(flweb.CountriesPath))
	s.Step(`^I should receive an "([^"]*)" response with the following countries:$`,
		assertResponseStatusAndListStep(&[]*flentities.Country{}))
}
