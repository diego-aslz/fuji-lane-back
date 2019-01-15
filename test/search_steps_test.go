package fujilane

import (
	"github.com/DATA-DOG/godog"
	"github.com/nerde/fuji-lane-back/flweb"
)

func SearchContext(s *godog.Suite) {
	s.Step(`^I search for units with the following filters:$`, performGETWithQueryStep(flweb.SearchPath))
}
