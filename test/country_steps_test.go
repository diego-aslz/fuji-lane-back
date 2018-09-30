package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/fuji-lane-back/flentities"
)

func theFollowingCountries(table *gherkin.DataTable) error {
	return createEntitiesFromTable(new(flentities.Country), table)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^the following countries:$`, theFollowingCountries)
}
