package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

func FeatureContext(s *godog.Suite) {
	s.Step(`^the following countries:$`, createFromTableStep(new(flentities.Country)))
	s.Step(`^we should have the following countries:$`, assertDatabaseRecordsStep(&[]*flentities.Country{}))
	s.Step(`^I list countries$`, performGETStep(flweb.CountriesPath))
	s.Step(`^the system should respond with "([^"]*)" and the following countries:$`,
		assertResponseStatusAndListStep(&[]*flentities.Country{}))
}
