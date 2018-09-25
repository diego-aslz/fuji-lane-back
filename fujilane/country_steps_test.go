package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func theFollowingCountries(table *gherkin.DataTable) error {
	return createFromTable(new(Country), table)
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^the following countries:$`, theFollowingCountries)
}
