package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/flweb"
)

func assertCountries(table *gherkin.DataTable) error {
	return assertDatabaseRecords(&[]*flentities.Country{}, table)
}

func requestCountries() error {
	return performGET(flweb.CountriesPath)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^the following countries:$`, createFromTableStep(new(flentities.Country)))
	s.Step(`^we should have the following countries:$`, assertCountries)
	s.Step(`^I list countries$`, requestCountries)
}
